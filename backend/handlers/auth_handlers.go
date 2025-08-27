package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"backend/entities"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db        *sql.DB
	jwtSecret []byte
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret_change_me"
	}
	return &AuthHandler{db: db, jwtSecret: []byte(secret)}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req entities.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	var id int
	err = h.db.QueryRow("INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", req.Login, string(hash)).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]any{"id": id, "login": req.Login})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req entities.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Login == "" || req.Password == "" {
		http.Error(w, "login and password are required", http.StatusBadRequest)
		return
	}

	var (
		id           int
		login        string
		passwordHash string
	)
	err := h.db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = $1", req.Login).Scan(&id, &login, &passwordHash)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"sub":   id,
		"login": login,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(h.jwtSecret)
	if err != nil {
		http.Error(w, "failed to sign token", http.StatusInternalServerError)
		return
	}

	resp := entities.AuthResponse{Token: signed}
	resp.User.ID = id
	resp.User.Login = login
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import { API_BASE_URL } from '../config';

export default function Register() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      await axios.post(`${API_BASE_URL}/auth/register`, { login, password });
      setSuccess('Регистрация успешна. Теперь вы можете войти.');
      setError(null);
      setTimeout(() => navigate('/login'), 800);
    } catch (err) {
      setError(err.response?.data || 'Ошибка регистрации');
      setSuccess(null);
    }
  };

  return (
    <div className="container">
      <h2>Регистрация</h2>
      {error && <div className="error">{String(error)}</div>}
      {success && <div className="success">{success}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Логин</label>
          <input value={login} onChange={(e) => setLogin(e.target.value)} />
        </div>
        <div className="form-group">
          <label>Пароль</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <button className="btn" type="submit">Зарегистрироваться</button>
      </form>
      <p>Уже есть аккаунт? <Link to="/login">Войти</Link></p>
    </div>
  );
} 
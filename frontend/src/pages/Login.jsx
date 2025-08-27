import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate, Link } from 'react-router-dom';
import { API_BASE_URL } from '../config';

export default function Login() {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const res = await axios.post(`${API_BASE_URL}/auth/login`, { login, password });
      localStorage.setItem('token', res.data.token);
      localStorage.setItem('user', JSON.stringify(res.data.user));
      setError(null);
      navigate('/');
    } catch (err) {
      setError(err.response?.data || 'Ошибка авторизации');
    }
  };

  return (
    <div className="container">
      <h2>Вход</h2>
      {error && <div className="error">{String(error)}</div>}
      <form onSubmit={handleSubmit}>
        <div className="form-group">
          <label>Логин</label>
          <input value={login} onChange={(e) => setLogin(e.target.value)} />
        </div>
        <div className="form-group">
          <label>Пароль</label>
          <input type="password" value={password} onChange={(e) => setPassword(e.target.value)} />
        </div>
        <button className="btn" type="submit">Войти</button>
      </form>
      <p>Нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link></p>
    </div>
  );
} 
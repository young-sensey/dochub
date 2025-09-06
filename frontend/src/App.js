import React from 'react';
import { Routes, Route, Link, Navigate, useLocation, useNavigate } from 'react-router-dom';
import './App.css';

import DocumentsList from './pages/DocumentsList';
import DocumentCreate from './pages/DocumentCreate';
import DocumentEdit from './pages/DocumentEdit';
import DocumentView from './pages/DocumentView';
import CategoriesList from './pages/CategoriesList';
import CategoryCreate from './pages/CategoryCreate';
import CategoryEdit from './pages/CategoryEdit';
import CategoryView from './pages/CategoryView';
import DocumentsByCategory from './pages/DocumentsByCategory';
import Sidebar from './components/Sidebar';
import Login from './pages/Login';
import Register from './pages/Register';

function Private({ children }) {
  const token = localStorage.getItem('token');
  const location = useLocation();
  if (!token) return <Navigate to="/login" replace state={{ from: location }} />;
  return children;
}

function HeaderAuth() {
  const navigate = useNavigate();
  const token = localStorage.getItem('token');
  const user = localStorage.getItem('user');
  const onLogout = () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    navigate('/login');
  };
  return (
    <nav>
      {token ? (
        <>
          <span className="user-login">{JSON.parse(user).login}</span>
          <span className="btn" onClick={onLogout}>Выйти</span>
        </>
      ) : ''}
    </nav>
  );
}

function App() {
  return (
    <div className="App">
      <div className="header">
        <div className="container">
          <p className="header-title">Система управления документами</p>
          <HeaderAuth />
        </div>
      </div>

      <div className="main-content">
        <Sidebar />
        <div className="content-area">
          <Routes>
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />

            <Route path="/" element={<Private><DocumentsList /></Private>} />
            <Route path="/documents/new" element={<Private><DocumentCreate /></Private>} />
            <Route path="/documents/:id" element={<Private><DocumentView /></Private>} />
            <Route path="/documents/:id/edit" element={<Private><DocumentEdit /></Private>} />
            <Route path="/categories" element={<Private><CategoriesList /></Private>} />
            <Route path="/categories/new" element={<Private><CategoryCreate /></Private>} />
            <Route path="/categories/:id" element={<Private><CategoryView /></Private>} />
            <Route path="/categories/:id/edit" element={<Private><CategoryEdit /></Private>} />
            <Route path="/category/:categoryId" element={<Private><DocumentsByCategory /></Private>} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </div>
    </div>
  );
}

export default App; 
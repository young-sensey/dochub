import React, { useEffect, useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config';

export default function Sidebar() {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const location = useLocation();
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;

  useEffect(() => {
    if (!token) {
      setLoading(false);
      setCategories([]);
      return;
    }
    const fetchCategories = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/categories`);
        setCategories(response.data);
      } catch (err) {
        console.error('Ошибка при загрузке категорий:', err);
      } finally {
        setLoading(false);
      }
    };
    fetchCategories();
  }, [token]);

  const isActive = (path) => {
    return location.pathname === path;
  };

  const isCategoryActive = (categoryId) => {
    return location.pathname === `/category/${categoryId}`;
  };

  const isNoCategoryActive = () => {
    return location.pathname === `/category/null`;
  };

  if (!token) {
    return null;
  }

  if (loading) {
    return (
      <div className="sidebar">
        <div className="sidebar-header">
          <h3>Категории</h3>
        </div>
        <div className="sidebar-content">
          <div className="loading">Загрузка...</div>
        </div>
      </div>
    );
  }

  return (
    <div className="sidebar">
      <div className="sidebar-content">
        <Link 
          to="/" 
          className={`sidebar-item ${isActive('/') ? 'active' : ''}`}
        >
          Все документы
        </Link>
        <Link 
          to="/category/null" 
          className={`sidebar-item ${isNoCategoryActive() ? 'active' : ''}`}
        >
          Без категории
        </Link>
        {categories.map((category) => (
          <Link
            key={category.id}
            to={`/category/${category.id}`}
            className={`sidebar-item ${isCategoryActive(category.id) ? 'active' : ''}`}
          >
            {category.name}
          </Link>
        ))}
        <Link
          to="/categories" 
          className={`sidebar-item categories ${isActive('/categories') ? 'active' : ''}`}
        >
          Категории
        </Link>
      </div>
    </div>
  );
} 
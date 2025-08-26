import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config';

export default function CategoryView() {
  const { id } = useParams();
  const [categoryData, setCategoryData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const loadCategory = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/categories/${id}`);
        setCategoryData(response.data);
      } catch (err) {
        setError('Ошибка при загрузке категории: ' + err.message);
      } finally {
        setLoading(false);
      }
    };
    loadCategory();
  }, [id]);

  if (loading) return <div><div className="loading">Загрузка...</div></div>;
  if (error) return <div><div className="error">{error}</div></div>;
  if (!categoryData) return <div><div className="error">Категория не найдена</div></div>;

  return (
    <div>
      <h2>{categoryData.name}</h2>
      <div style={{ whiteSpace: 'pre-wrap', marginBottom: 16 }}>{categoryData.description}</div>
      <div style={{ color: '#666', fontSize: 14, marginBottom: 16 }}>
        Создана: {new Date(categoryData.created_at).toLocaleString('ru-RU')} | 
        Обновлена: {new Date(categoryData.updated_at).toLocaleString('ru-RU')}
      </div>

      <div style={{ display: 'flex', gap: 8 }}>
        <Link to={`/categories/${id}/edit`} className="btn btn-warning">Редактировать</Link>
        <Link to="/categories" className="btn">Назад к списку</Link>
      </div>
    </div>
  );
} 
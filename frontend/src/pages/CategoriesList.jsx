import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config';

export default function CategoriesList() {
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const fetchCategories = async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_BASE_URL}/categories`);
      setCategories(response.data);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке категорий: ' + err.message);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchCategories();
  }, []);

  useEffect(() => {
    if (success || error) {
      const timer = setTimeout(() => {
        setSuccess(null);
        setError(null);
      }, 4000);
      return () => clearTimeout(timer);
    }
  }, [success, error]);

  const deleteCategory = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить эту категорию?')) return;
    try {
      await axios.delete(`${API_BASE_URL}/categories/${id}`);
      setCategories((prev) => prev.filter((cat) => cat.id !== id));
      setSuccess('Категория успешно удалена!');
      setError(null);
    } catch (err) {
      setError('Ошибка при удалении категории: ' + err.message);
    }
  };

  return (
    <div>
      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}

      <div className="documents-list">
        <div className="documents-list-header">
          <h2>Список категорий</h2>
          <Link to="/categories/new" className="btn btn-primary">Добавить категорию</Link>
        </div>

        {loading ? (
          <div className="loading">Загрузка категорий...</div>
        ) : categories.length === 0 ? (
          <div className="loading">Категории не найдены</div>
        ) : (
          categories.map((category) => (
            <div key={category.id} className="document-item">
              <div className="document-info">
                <div className="document-title">{category.name}</div>
                <div className="document-content">{category.description}</div>
                <div className="document-meta">
                  Создана: {new Date(category.created_at).toLocaleString('ru-RU')} | 
                  Обновлена: {new Date(category.updated_at).toLocaleString('ru-RU')}
                </div>
              </div>
              <div className="document-actions">
                <Link className="btn" to={`/categories/${category.id}`}>Просмотр</Link>
                <Link className="btn btn-warning" to={`/categories/${category.id}/edit`}>Редактировать</Link>
                <button className="btn btn-danger" onClick={() => deleteCategory(category.id)}>Удалить</button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
} 
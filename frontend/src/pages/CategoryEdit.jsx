import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import axios from 'axios';
import CategoryForm from '../components/CategoryForm';
import { API_BASE_URL } from '../config';

export default function CategoryEdit() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [categoryValues, setCategoryValues] = useState({ name: '', description: '' });

  useEffect(() => {
    const loadCategory = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/categories/${id}`);
        const data = response.data;
        setCategoryValues({ name: data.name, description: data.description });
      } catch (err) {
        setError('Ошибка при загрузке категории: ' + err.message);
      } finally {
        setLoading(false);
      }
    };
    loadCategory();
  }, [id]);

  const handleUpdate = async (values) => {
    try {
      await axios.put(`${API_BASE_URL}/categories/${id}`, values);
      navigate('/categories');
    } catch (err) {
      setError('Ошибка при обновлении категории: ' + err.message);
    }
  };

  if (loading) return <div><div className="loading">Загрузка...</div></div>;

  return (
    <div>
      <h2>Редактировать категорию</h2>
      {error && <div className="error">{error}</div>}

      <CategoryForm
        initialValues={categoryValues}
        submitButtonLabel="Обновить"
        onSubmit={handleUpdate}
      />
    </div>
  );
} 
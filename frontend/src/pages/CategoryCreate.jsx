import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import CategoryForm from '../components/CategoryForm';
import { API_BASE_URL } from '../config';

export default function CategoryCreate() {
  const navigate = useNavigate();
  const [error, setError] = useState(null);

  const handleCreate = async (values) => {
    try {
      await axios.post(`${API_BASE_URL}/categories`, values);
      navigate('/categories');
    } catch (err) {
      setError('Ошибка при создании категории: ' + err.message);
    }
  };

  return (
    <div>
      <h2>Создать новую категорию</h2>
      {error && <div className="error">{error}</div>}

      <CategoryForm
        initialValues={{ name: '', description: '' }}
        submitButtonLabel="Создать"
        onSubmit={handleCreate}
      />
    </div>
  );
}
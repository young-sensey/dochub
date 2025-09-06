import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';
import DocumentForm from '../components/DocumentForm';
import { API_BASE_URL } from '../config';

export default function DocumentCreate() {
  const navigate = useNavigate();
  const [error, setError] = useState(null);

  const handleCreate = async (values, file) => {
    try {
      const formData = new FormData();
      formData.append('title', values.title);
      formData.append('content', values.content);
      if (file) formData.append('file', file);

      await axios.post(`${API_BASE_URL}/dock`, formData, {
        headers: { 'Content-Type': 'multipart/form-data' },
      });

      navigate('/');
    } catch (err) {
      setError('Ошибка при создании документа: ' + err.message);
    }
  };

  return (
    <div>
      <h2>Создать новый документ</h2>
      {error && <div className="error">{error}</div>}

      <DocumentForm
        initialValues={{ title: '', content: '', category_id: null }}
        allowFileUpload={true}
        submitButtonLabel="Создать"
        onSubmit={handleCreate}
      />


    </div>
  );
} 
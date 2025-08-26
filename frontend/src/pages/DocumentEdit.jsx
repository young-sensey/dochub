import React, { useEffect, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import axios from 'axios';
import DocumentForm from '../components/DocumentForm';
import { API_BASE_URL } from '../config';

export default function DocumentEdit() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [error, setError] = useState(null);
  const [loading, setLoading] = useState(true);
  const [documentValues, setDocumentValues] = useState({ title: '', content: '', author: '' });

  useEffect(() => {
    const loadDocument = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/dock/${id}`);
        const data = response.data;
        setDocumentValues({ 
          title: data.title, 
          content: data.content, 
          author: data.author, 
          category_id: data.category_id 
        });
      } catch (err) {
        setError('Ошибка при загрузке документа: ' + err.message);
      } finally {
        setLoading(false);
      }
    };
    loadDocument();
  }, [id]);

  const handleUpdate = async (values) => {
    try {
      await axios.put(`${API_BASE_URL}/dock/${id}`, values);
      navigate('/');
    } catch (err) {
      setError('Ошибка при обновлении документа: ' + err.message);
    }
  };

  if (loading) return <div><div className="loading">Загрузка...</div></div>;

  return (
    <div>
      <h2>Редактировать документ</h2>
      {error && <div className="error">{error}</div>}

      <DocumentForm
        initialValues={documentValues}
        allowFileUpload={false}
        submitButtonLabel="Обновить"
        onSubmit={handleUpdate}
      />


    </div>
  );
} 
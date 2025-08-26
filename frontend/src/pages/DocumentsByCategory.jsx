import React, { useEffect, useState, useCallback } from 'react';
import { Link, useParams } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config';

const DownloadIcon = () => (
  <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" fill="currentColor" viewBox="0 0 20 20">
    <path d="M10 2a1 1 0 0 1 1 1v8.586l2.293-2.293a1 1 0 1 1 1.414 1.414l-4 4a1 1 0 0 1-1.414 0l-4-4A1 1 0 1 1 5.707 9.293L8 11.586V3a1 1 0 0 1 1-1zm-7 13a1 1 0 0 1 1 1v1h12v-1a1 1 0 1 1 2 0v2a1 1 0 0 1-1 1H3a1 1 0 0 1-1-1v-2a1 1 0 0 1 1-1z"/>
  </svg>
);

export default function DocumentsByCategory() {
  const { categoryId } = useParams();
  const [documents, setDocuments] = useState([]);
  const [category, setCategory] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [success, setSuccess] = useState(null);

  const fetchDocuments = useCallback(async () => {
    try {
      setLoading(true);
      const response = await axios.get(`${API_BASE_URL}/dock?category_id=${categoryId}`);
      setDocuments(response.data || []);
      setError(null);
    } catch (err) {
      setError('Ошибка при загрузке документов: ' + err.message);
      setDocuments([]);
    } finally {
      setLoading(false);
    }
  }, [categoryId]);

  const fetchCategory = useCallback(async () => {
    try {
      if (categoryId === 'null') {
        setCategory({ name: 'Без категории', description: 'Документы без назначенной категории' });
      } else {
        const response = await axios.get(`${API_BASE_URL}/categories/${categoryId}`);
        setCategory(response.data);
      }
    } catch (err) {
      console.error('Ошибка при загрузке категории:', err);
    }
  }, [categoryId]);

  useEffect(() => {
    fetchCategory();
    fetchDocuments();
  }, [fetchCategory, fetchDocuments]);

  useEffect(() => {
    if (success || error) {
      const timer = setTimeout(() => {
        setSuccess(null);
        setError(null);
      }, 4000);
      return () => clearTimeout(timer);
    }
  }, [success, error]);

  const handleDownload = async (id, filePath) => {
    try {
      const response = await axios.get(`/dock/${id}/download`, {
        responseType: 'blob'
      });
      
      // Получаем имя файла из пути
      const fileName = filePath ? filePath.split('/').pop() : `document_${id}`;
      
      // Создаем ссылку для скачивания
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', fileName);
      document.body.appendChild(link);
      link.click();
      link.remove();
      window.URL.revokeObjectURL(url);
    } catch (err) {
      setError('Ошибка при скачивании файла: ' + err.message);
    }
  };

  const deleteDocument = async (id) => {
    if (!window.confirm('Вы уверены, что хотите удалить этот документ?')) return;
    try {
      await axios.delete(`${API_BASE_URL}/dock/${id}`);
      setDocuments((prev) => prev.filter((doc) => doc.id !== id));
      setSuccess('Документ успешно удален!');
      setError(null);
    } catch (err) {
      setError('Ошибка при удалении документа: ' + err.message);
    }
  };

  if (loading) return <div><div className="loading">Загрузка...</div></div>;
  if (error) return <div><div className="error">{error}</div></div>;
  if (!category) return <div><div className="error">Категория не найдена</div></div>;

  return (
    <div>
      {error && <div className="error">{error}</div>}
      {success && <div className="success">{success}</div>}

      <div className="documents-list" style={{ marginTop: 16 }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <h2>Документы категории: {category.name}</h2>
          <Link to="/documents/new" className="btn btn-primary">Добавить документ</Link>
        </div>

        {!documents || documents.length === 0 ? (
          <div className="loading">В этой категории нет документов</div>
        ) : (
          documents.map((doc) => (
            <div key={doc.id} className="document-item">
              <div className="document-info">
                <div className="document-title">{doc.title}</div>
                <div className="document-content">{doc.content}</div>
                <div className="document-meta">
                  Автор: {doc.author} | 
                  Создан: {new Date(doc.created_at).toLocaleString('ru-RU')} | 
                  Обновлен: {new Date(doc.updated_at).toLocaleString('ru-RU')}
                </div>
              </div>
              <div className="document-actions">
                {doc.file_path && (
                  <button
                    className="btn btn-success"
                    onClick={() => handleDownload(doc.id, doc.file_path)}
                    title="Скачать файл"
                    style={{ display: 'inline-flex', alignItems: 'center', justifyContent: 'center', padding: '10px 16px' }}
                  >
                    <DownloadIcon />
                  </button>
                )}
                <Link className="btn" to={`/documents/${doc.id}`}>Просмотр</Link>
                <Link className="btn btn-warning" to={`/documents/${doc.id}/edit`}>Редактировать</Link>
                <button className="btn btn-danger" onClick={() => deleteDocument(doc.id)}>Удалить</button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
} 
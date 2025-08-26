import React, { useEffect, useState } from 'react';
import { Link, useParams } from 'react-router-dom';
import axios from 'axios';
import { API_BASE_URL } from '../config';

export default function DocumentView() {
  const { id } = useParams();
  const [documentData, setDocumentData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const loadDocument = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/dock/${id}`);
        setDocumentData(response.data);
      } catch (err) {
        setError('Ошибка при загрузке документа: ' + err.message);
      } finally {
        setLoading(false);
      }
    };
    loadDocument();
  }, [id]);

  const handleDownload = async () => {
    try {
      const response = await axios.get(`/dock/${id}/download`, {
        responseType: 'blob'
      });
      
      // Получаем имя файла из пути
      const fileName = documentData.file_path ? documentData.file_path.split('/').pop() : `document_${id}`;
      
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

  if (loading) return <div><div className="loading">Загрузка...</div></div>;
  if (error) return <div><div className="error">{error}</div></div>;
  if (!documentData) return <div><div className="error">Документ не найден</div></div>;

  return (
    <div>
      <h2>{documentData.title}</h2>
      <div style={{ marginBottom: 8 }}>Автор: {documentData.author}</div>
      <div style={{ whiteSpace: 'pre-wrap', marginBottom: 16 }}>{documentData.content}</div>
      <div style={{ color: '#666', fontSize: 14, marginBottom: 16 }}>
        Создан: {new Date(documentData.created_at).toLocaleString('ru-RU')} | Обновлен: {new Date(documentData.updated_at).toLocaleString('ru-RU')}
      </div>

      <div style={{ display: 'flex', gap: 8 }}>
        <Link to={`/documents/${id}/edit`} className="btn btn-warning">Редактировать</Link>
        {documentData.file_path && (
          <button className="btn btn-success" onClick={handleDownload}>
            Скачать файл
          </button>
        )}
        <Link to="/" className="btn">Назад к списку</Link>
      </div>
    </div>
  );
} 
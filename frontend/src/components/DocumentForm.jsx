import React, { useState, useEffect } from 'react';
import axios from 'axios';
import { API_BASE_URL } from '../config';

export default function DocumentForm({
  initialValues = { title: '', content: '', author: '', category_id: null },
  allowFileUpload = false,
  submitButtonLabel = 'Сохранить',
  onSubmit,
}) {
  const [values, setValues] = useState(initialValues);
  const [file, setFile] = useState(null);
  const [categories, setCategories] = useState([]);
  const [loadingCategories, setLoadingCategories] = useState(true);

  useEffect(() => {
    setValues(initialValues);
  }, [initialValues]);

  useEffect(() => {
    const fetchCategories = async () => {
      try {
        const response = await axios.get(`${API_BASE_URL}/categories`);
        setCategories(response.data);
      } catch (err) {
        console.error('Ошибка при загрузке категорий:', err);
      } finally {
        setLoadingCategories(false);
      }
    };
    fetchCategories();
  }, []);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setValues((prev) => ({ 
      ...prev, 
      [name]: name === 'category_id' ? (value === '' ? null : parseInt(value)) : value 
    }));
  };

  const handleFileChange = (event) => {
    setFile(event.target.files && event.target.files[0] ? event.target.files[0] : null);
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    onSubmit(values, file);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="form-group">
        <label htmlFor="title">Заголовок:</label>
        <input
          type="text"
          id="title"
          name="title"
          value={values.title}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="content">Содержимое:</label>
        <textarea
          id="content"
          name="content"
          value={values.content}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="author">Автор:</label>
        <input
          type="text"
          id="author"
          name="author"
          value={values.author}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="category_id">Категория:</label>
        <select
          id="category_id"
          name="category_id"
          value={values.category_id || ''}
          onChange={handleChange}
        >
          <option value="">Без категории</option>
          {!loadingCategories && categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      {allowFileUpload && (
        <div className="form-group">
          <label htmlFor="file">Файл:</label>
          <input type="file" id="file" name="file" onChange={handleFileChange} />
        </div>
      )}

      <div>
        <button type="submit" className="btn btn-primary">{submitButtonLabel}</button>
      </div>
    </form>
  );
} 
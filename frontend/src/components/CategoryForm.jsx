import React, { useState, useEffect } from 'react';

export default function CategoryForm({
  initialValues = { name: '', description: '' },
  submitButtonLabel = 'Сохранить',
  onSubmit,
}) {
  const [values, setValues] = useState(initialValues);

  useEffect(() => {
    setValues(initialValues);
  }, [initialValues]);

  const handleChange = (event) => {
    const { name, value } = event.target;
    setValues((prev) => ({ ...prev, [name]: value }));
  };

  const handleSubmit = (event) => {
    event.preventDefault();
    onSubmit(values);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="form-group">
        <label htmlFor="name">Название:</label>
        <input
          type="text"
          id="name"
          name="name"
          value={values.name}
          onChange={handleChange}
          required
        />
      </div>

      <div className="form-group">
        <label htmlFor="description">Описание:</label>
        <textarea
          id="description"
          name="description"
          value={values.description}
          onChange={handleChange}
          rows="4"
        />
      </div>

      <div>
        <button type="submit" className="btn btn-primary">{submitButtonLabel}</button>
      </div>
    </form>
  );
} 
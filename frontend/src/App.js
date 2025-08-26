import React from 'react';
import { Routes, Route, Link, Navigate } from 'react-router-dom';
import './App.css';

import DocumentsList from './pages/DocumentsList';
import DocumentCreate from './pages/DocumentCreate';
import DocumentEdit from './pages/DocumentEdit';
import DocumentView from './pages/DocumentView';
import CategoriesList from './pages/CategoriesList';
import CategoryCreate from './pages/CategoryCreate';
import CategoryEdit from './pages/CategoryEdit';
import CategoryView from './pages/CategoryView';
import DocumentsByCategory from './pages/DocumentsByCategory';
import Sidebar from './components/Sidebar';

function App() {
  return (
    <div className="App">
      <div className="header">
        <div className="container" style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <h1 style={{ margin: 0 }}>Система управления документами</h1>
          <nav>
            <Link to="/" className="btn">Документы</Link>
            <Link to="/categories" className="btn">Категории</Link>
          </nav>
        </div>
      </div>

      <div className="main-content">
        <Sidebar />
        <div className="content-area">
          <Routes>
            <Route path="/" element={<DocumentsList />} />
            <Route path="/documents/new" element={<DocumentCreate />} />
            <Route path="/documents/:id" element={<DocumentView />} />
            <Route path="/documents/:id/edit" element={<DocumentEdit />} />
            <Route path="/categories" element={<CategoriesList />} />
            <Route path="/categories/new" element={<CategoryCreate />} />
            <Route path="/categories/:id" element={<CategoryView />} />
            <Route path="/categories/:id/edit" element={<CategoryEdit />} />
            <Route path="/category/:categoryId" element={<DocumentsByCategory />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </div>
    </div>
  );
}

export default App; 
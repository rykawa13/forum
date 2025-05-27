// src/utils/AdminRoute.js
import React from 'react';
import { useSelector } from 'react-redux';
import { Navigate } from 'react-router-dom';

const AdminRoute = ({ children }) => {
  const { user, isAuthenticated } = useSelector(state => state.auth);
  
  // Проверяем, что пользователь авторизован И является администратором
  if (!isAuthenticated || !user?.is_admin) {
    return <Navigate to="/" replace />;
  }

  return children;
};

export default AdminRoute;
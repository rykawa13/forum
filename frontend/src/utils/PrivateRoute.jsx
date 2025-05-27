import React from 'react';
import { useSelector } from 'react-redux';
import { Navigate, useLocation } from 'react-router-dom';

/**
 * Компонент для защиты маршрутов, требующих аутентификации
 * @param children - Защищаемый компонент
 * @param redirectPath - Путь для перенаправления (по умолчанию '/login')
 */
const PrivateRoute = ({ children, redirectPath = '/login' }) => {
  const { isAuthenticated } = useSelector(state => state.auth);
  const location = useLocation();

  if (!isAuthenticated) {
    // Сохраняем текущий путь для последующего перенаправления после входа
    return <Navigate to={redirectPath} state={{ from: location }} replace />;
  }

  return children;
};

export default PrivateRoute;
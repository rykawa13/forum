import React from 'react';
import { Container, Typography, Box } from '@mui/material';
import UserList from './UserList';
import { useAuth } from '../../hooks/useAuth';
import { Navigate } from 'react-router-dom';

const AdminPanel: React.FC = () => {
  const { user } = useAuth();

  if (!user?.is_admin) {
    return <Navigate to="/" replace />;
  }

  return (
    <Container>
      <Box sx={{ py: 4 }}>
        <Typography variant="h3" component="h1" gutterBottom>
          Панель администратора
        </Typography>
        <UserList />
      </Box>
    </Container>
  );
};

export default AdminPanel; 
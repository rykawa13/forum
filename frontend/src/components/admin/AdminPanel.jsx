import React from 'react';
import { Box, Typography, Container, Divider } from '@mui/material';
import { useSelector } from 'react-redux';
import UserList from './UserList';
import ForumStats from './ForumStats';

const AdminPanel = () => {
  const { user } = useSelector(state => state.auth);

  return (
    <Container maxWidth="lg">
      <Box sx={{ mt: 4, mb: 4 }}>
        <Typography variant="h4" component="h1" gutterBottom>
          Панель администратора
        </Typography>
        <Typography variant="subtitle1" color="text.secondary" gutterBottom>
          Добро пожаловать, {user?.username}!
        </Typography>

        <Box sx={{ mt: 4, mb: 4 }}>
          <Typography variant="h5" gutterBottom>
            Статистика форума
          </Typography>
          <ForumStats />
        </Box>

        <Divider sx={{ my: 4 }} />

        <Box sx={{ mt: 4 }}>
          <Typography variant="h5" gutterBottom>
            Управление пользователями
          </Typography>
          <UserList />
        </Box>
      </Box>
    </Container>
  );
};

export default AdminPanel; 
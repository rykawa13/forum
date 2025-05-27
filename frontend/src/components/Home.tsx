import React from 'react';
import { Container, Typography, Box } from '@mui/material';
import { useAuth } from '../hooks/useAuth';

const Home: React.FC = () => {
  const { user } = useAuth();

  return (
    <Container>
      <Box sx={{ mt: 4 }}>
        <Typography variant="h4" gutterBottom>
          Добро пожаловать, {user?.username}!
        </Typography>
        <Typography variant="body1">
          {user?.is_admin 
            ? 'У вас есть права администратора. Используйте иконку в правом верхнем углу для доступа к панели управления.'
            : 'Вы вошли как обычный пользователь.'}
        </Typography>
      </Box>
    </Container>
  );
};

export default Home; 
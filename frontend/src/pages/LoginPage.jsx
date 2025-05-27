import React from 'react';
import { Container, Box, Typography } from '@mui/material';
import LoginForm from '../components/auth/LoginForm';
import { Link } from 'react-router-dom';

const LoginPage = () => {
  return (
    <Container maxWidth="xs">
      <Box sx={{ mt: 8, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Typography component="h1" variant="h5">
          Вход в систему
        </Typography>
        <LoginForm />
        <Box sx={{ mt: 2 }}>
          <Typography variant="body2">
            Нет аккаунта? <Link to="/register">Зарегистрируйтесь</Link>
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default LoginPage;
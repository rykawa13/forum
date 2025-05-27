import React from 'react';
import { Container, Box, Typography, Link } from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import RegisterForm from '../components/auth/RegisterForm';

const RegisterPage = () => {
  return (
    <Container maxWidth="xs">
      <Box sx={{
        mt: 8,
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
      }}>
        <Typography component="h1" variant="h5" sx={{ mb: 2 }}>
          Регистрация
        </Typography>
        <RegisterForm />
        <Box sx={{ mt: 3 }}>
          <Typography variant="body2">
            Уже есть аккаунт?{' '}
            <Link component={RouterLink} to="/login">
              Войдите
            </Link>
          </Typography>
        </Box>
      </Box>
    </Container>
  );
};

export default RegisterPage;
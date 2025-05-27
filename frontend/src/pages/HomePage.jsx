import React from 'react';
import { Box, Typography, Button, Container } from '@mui/material';
import { Link } from 'react-router-dom';
import { useSelector } from 'react-redux';

const HomePage = () => {
  const { isAuthenticated } = useSelector(state => state.auth);

  return (
    <Container maxWidth="lg">
      <Box sx={{ 
        textAlign: 'center', 
        py: 10,
        background: 'linear-gradient(rgba(255,255,255,0.9), rgba(255,255,255,0.9))',
        borderRadius: 4,
        mt: 4,
        p: 4
      }}>
        <Typography variant="h2" component="h1" gutterBottom sx={{ fontWeight: 'bold' }}>
          Добро пожаловать на наш форум
        </Typography>
        <Typography variant="h5" component="p" gutterBottom sx={{ mb: 4 }}>
          Общайтесь, делитесь знаниями и находите единомышленников
        </Typography>
        {!isAuthenticated ? (
          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
            <Button 
              variant="contained" 
              size="large" 
              component={Link} 
              to="/register"
              sx={{ px: 4 }}
            >
              Регистрация
            </Button>
            <Button 
              variant="outlined" 
              size="large" 
              component={Link} 
              to="/login"
              sx={{ px: 4 }}
            >
              Вход
            </Button>
          </Box>
        ) : (
          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center' }}>
            <Button 
              variant="contained" 
              size="large" 
              component={Link} 
              to="/forum"
              sx={{ px: 4 }}
            >
              На форум
            </Button>
            <Button 
              variant="outlined" 
              size="large" 
              component={Link} 
              to="/chat"
              sx={{ px: 4 }}
            >
              В чат
            </Button>
          </Box>
        )}
      </Box>
    </Container>
  );
};

export default HomePage;
import React from 'react';
import { Container, Box, Typography, Avatar, Button } from '@mui/material';
import { useSelector } from 'react-redux';

const ProfilePage = () => {
  const { user } = useSelector(state => state.auth);

  return (
    <Container maxWidth="md">
      <Box sx={{ mt: 4, display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <Avatar sx={{ width: 100, height: 100, mb: 2 }}>
          {user?.username.charAt(0)}
        </Avatar>
        <Typography variant="h4" gutterBottom>
          {user?.username}
        </Typography>
        <Typography variant="subtitle1" gutterBottom>
          {user?.email}
        </Typography>
        {user?.isAdmin && (
          <Typography color="primary" sx={{ mt: 1 }}>
            Администратор
          </Typography>
        )}
        <Button variant="outlined" sx={{ mt: 3 }}>
          Редактировать профиль
        </Button>
      </Box>
    </Container>
  );
};

export default ProfilePage;
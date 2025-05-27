import React from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  IconButton,
} from '@mui/material';
import { Link as RouterLink } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { AdminPanelSettings } from '@mui/icons-material';

const Header: React.FC = () => {
  const { user, logout } = useAuth();

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component={RouterLink} to="/" sx={{ flexGrow: 1, textDecoration: 'none', color: 'inherit' }}>
          Auth Service
        </Typography>
        <Box>
          {user ? (
            <>
              {user.is_admin && (
                <IconButton
                  component={RouterLink}
                  to="/admin"
                  color="inherit"
                  sx={{ mr: 2 }}
                >
                  <AdminPanelSettings />
                </IconButton>
              )}
              <Button color="inherit" onClick={logout}>
                Выйти
              </Button>
            </>
          ) : (
            <>
              <Button color="inherit" component={RouterLink} to="/login">
                Войти
              </Button>
              <Button color="inherit" component={RouterLink} to="/register">
                Регистрация
              </Button>
            </>
          )}
        </Box>
      </Toolbar>
    </AppBar>
  );
};

export default Header; 
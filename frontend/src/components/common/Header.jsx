import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useSelector, useDispatch } from 'react-redux';
import { logout } from '../../store/authSlice';
import { AppBar, Toolbar, Typography, Button, Box, IconButton } from '@mui/material';
import AdminPanelSettingsIcon from '@mui/icons-material/AdminPanelSettings';

const Header = () => {
  const { isAuthenticated, user } = useSelector(state => state.auth);
  const dispatch = useDispatch();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await dispatch(logout());
    navigate('/');
  };

  return (
    <AppBar position="static">
      <Toolbar>
        <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
          <Link to="/" style={{ color: 'white', textDecoration: 'none' }}>
            Форум
          </Link>
        </Typography>
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
          <Button color="inherit" component={Link} to="/forum">
            ФОРУМ
          </Button>
          <Button color="inherit" component={Link} to="/chat">
            ЧАТ
          </Button>
          {isAuthenticated ? (
            <>
              {user?.is_admin && (
                <IconButton
                  color="inherit"
                  component={Link}
                  to="/admin/"
                  sx={{ ml: 1 }}
                  title="Админ-панель"
                >
                  <AdminPanelSettingsIcon />
                </IconButton>
              )}
              <Button color="inherit" component={Link} to="/profile">
                {user?.username}
              </Button>
              <Button color="inherit" onClick={handleLogout}>
                ВЫЙТИ
              </Button>
            </>
          ) : (
            <>
              <Button color="inherit" component={Link} to="/login">
                Войти
              </Button>
              <Button color="inherit" component={Link} to="/register">
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
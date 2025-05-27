import React from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { TextField, Button, Box, Typography } from '@mui/material';
import { login } from '../../store/authSlice';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate, useLocation } from 'react-router-dom';
import { useEffect } from 'react';

const LoginForm = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const location = useLocation();
  const { isAuthenticated, error, loading } = useSelector(state => state.auth);

  useEffect(() => {
    if (isAuthenticated) {
      // Перенаправляем на предыдущую страницу или на форум
      const from = location.state?.from?.pathname || '/forum';
      navigate(from, { replace: true });
    }
  }, [isAuthenticated, navigate, location]);

  const formik = useFormik({
    initialValues: {
      email: '',
      password: '',
    },
    validationSchema: Yup.object({
      email: Yup.string().email('Некорректный email').required('Обязательное поле'),
      password: Yup.string().required('Обязательное поле'),
    }),
    onSubmit: async (values) => {
      dispatch(login(values));
    },
  });

  // Получаем текст ошибки из объекта error
  const errorMessage = error?.error || error;

  return (
    <Box component="form" onSubmit={formik.handleSubmit} sx={{ mt: 3 }}>
      {errorMessage && (
        <Typography color="error" sx={{ mb: 2 }}>
          {errorMessage}
        </Typography>
      )}
      <TextField
        fullWidth
        id="email"
        name="email"
        label="Email"
        value={formik.values.email}
        onChange={formik.handleChange}
        error={formik.touched.email && Boolean(formik.errors.email)}
        helperText={formik.touched.email && formik.errors.email}
        sx={{ mb: 2 }}
        disabled={loading}
      />
      <TextField
        fullWidth
        id="password"
        name="password"
        label="Пароль"
        type="password"
        value={formik.values.password}
        onChange={formik.handleChange}
        error={formik.touched.password && Boolean(formik.errors.password)}
        helperText={formik.touched.password && formik.errors.password}
        sx={{ mb: 2 }}
        disabled={loading}
      />
      <Button 
        type="submit" 
        fullWidth 
        variant="contained" 
        sx={{ mt: 2 }}
        disabled={loading}
      >
        {loading ? 'Вход...' : 'Войти'}
      </Button>
    </Box>
  );
};

export default LoginForm;
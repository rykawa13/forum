import React from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { TextField, Button, Box, Typography } from '@mui/material';
import { register } from '../../store/authSlice';
import { useDispatch, useSelector } from 'react-redux';
import { useNavigate } from 'react-router-dom';
import { useEffect } from 'react';

const RegisterForm = () => {
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const { isAuthenticated, error, loading } = useSelector(state => state.auth);

  useEffect(() => {
    if (isAuthenticated) {
      navigate('/forum', { replace: true });
    }
  }, [isAuthenticated, navigate]);

  const formik = useFormik({
    initialValues: {
      username: '',
      email: '',
      password: '',
      confirmPassword: '',
    },
    validationSchema: Yup.object({
      username: Yup.string()
        .min(3, 'Минимум 3 символа')
        .required('Обязательное поле'),
      email: Yup.string()
        .email('Некорректный email')
        .required('Обязательное поле'),
      password: Yup.string()
        .min(6, 'Минимум 6 символов')
        .required('Обязательное поле'),
      confirmPassword: Yup.string()
        .oneOf([Yup.ref('password'), null], 'Пароли должны совпадать')
        .required('Обязательное поле'),
    }),
    onSubmit: (values) => {
      dispatch(register({
        username: values.username,
        email: values.email,
        password: values.password,
      }));
    },
  });

  // Получаем текст ошибки из объекта error
  const errorMessage = error?.error || error;

  return (
    <Box component="form" onSubmit={formik.handleSubmit} sx={{ mt: 3 }}>
      {errorMessage && (
        <Typography color="error" sx={{ mb: 2 }}>
          {typeof errorMessage === 'string' ? errorMessage : 'Произошла ошибка при регистрации'}
        </Typography>
      )}
      <TextField
        fullWidth
        id="username"
        name="username"
        label="Имя пользователя"
        value={formik.values.username}
        onChange={formik.handleChange}
        error={formik.touched.username && Boolean(formik.errors.username)}
        helperText={formik.touched.username && formik.errors.username}
        sx={{ mb: 2 }}
        disabled={loading}
      />
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
      <TextField
        fullWidth
        id="confirmPassword"
        name="confirmPassword"
        label="Подтвердите пароль"
        type="password"
        value={formik.values.confirmPassword}
        onChange={formik.handleChange}
        error={formik.touched.confirmPassword && Boolean(formik.errors.confirmPassword)}
        helperText={formik.touched.confirmPassword && formik.errors.confirmPassword}
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
        {loading ? 'Регистрация...' : 'Зарегистрироваться'}
      </Button>
    </Box>
  );
};

export default RegisterForm;
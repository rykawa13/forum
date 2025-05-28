import React from 'react';
import { useSelector } from 'react-redux';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { TextField, Button, Box, Typography } from '@mui/material';

const PostForm = ({ onSubmit }) => {
  const { user, isAuthenticated } = useSelector(state => state.auth);

  const formik = useFormik({
    initialValues: {
      title: '',
      content: '',
    },
    validateOnMount: false,
    validateOnChange: true,
    validationSchema: Yup.object({
      title: Yup.string()
        .min(3, 'Минимум 3 символа')
        .max(100, 'Максимум 100 символов')
        .required('Обязательное поле'),
      content: Yup.string()
        .min(10, 'Минимум 10 символов')
        .required('Обязательное поле'),
    }),
    onSubmit: (values, { resetForm }) => {
      if (!isAuthenticated || !user?.id) {
        alert('Пожалуйста, войдите в систему, чтобы создать пост');
        return;
      }
      onSubmit({ ...values, author_id: user.id });
      resetForm();
    },
  });

  if (!isAuthenticated) {
    return (
      <Box sx={{ mb: 4 }}>
        <Typography color="error" sx={{ mb: 2 }}>
          Для создания поста необходимо войти в систему
        </Typography>
      </Box>
    );
  }

  const isFormValid = formik.dirty && !Object.keys(formik.errors).length;

  return (
    <Box component="form" onSubmit={formik.handleSubmit} sx={{ mb: 4 }}>
      <TextField
        fullWidth
        id="title"
        name="title"
        label="Заголовок"
        value={formik.values.title}
        onChange={formik.handleChange}
        onBlur={formik.handleBlur}
        error={formik.touched.title && Boolean(formik.errors.title)}
        helperText={formik.touched.title && formik.errors.title}
        sx={{ mb: 2 }}
      />
      <TextField
        fullWidth
        id="content"
        name="content"
        label="Содержание"
        multiline
        rows={4}
        value={formik.values.content}
        onChange={formik.handleChange}
        onBlur={formik.handleBlur}
        error={formik.touched.content && Boolean(formik.errors.content)}
        helperText={formik.touched.content && formik.errors.content}
        sx={{ mb: 2 }}
      />
      <Button 
        type="submit" 
        variant="contained" 
        disabled={formik.isSubmitting || !isFormValid}
      >
        Опубликовать
      </Button>
    </Box>
  );
};

export default PostForm;
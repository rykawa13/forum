import React from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { TextField, Button, Box } from '@mui/material';

const PostForm = ({ onSubmit }) => {
  const formik = useFormik({
    initialValues: {
      title: '',
      content: '',
    },
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
      onSubmit(values);
      resetForm();
    },
  });

  return (
    <Box component="form" onSubmit={formik.handleSubmit} sx={{ mb: 4 }}>
      <TextField
        fullWidth
        id="title"
        name="title"
        label="Заголовок"
        value={formik.values.title}
        onChange={formik.handleChange}
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
        error={formik.touched.content && Boolean(formik.errors.content)}
        helperText={formik.touched.content && formik.errors.content}
        sx={{ mb: 2 }}
      />
      <Button 
        type="submit" 
        variant="contained" 
        disabled={!formik.isValid || formik.isSubmitting}
      >
        Опубликовать
      </Button>
    </Box>
  );
};

export default PostForm;
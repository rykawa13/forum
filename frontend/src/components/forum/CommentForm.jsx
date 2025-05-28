import React from 'react';
import { useFormik } from 'formik';
import * as Yup from 'yup';
import { TextField, Button, Box } from '@mui/material';

const CommentForm = ({ postId, onSubmit }) => {
  const formik = useFormik({
    initialValues: {
      content: '',
    },
    validationSchema: Yup.object({
      content: Yup.string()
        .min(1, 'Минимум 1 символ')
        .required('Обязательное поле'),
    }),
    onSubmit: (values, { resetForm }) => {
      onSubmit({ ...values, postId });
      resetForm();
    },
  });

  return (
    <Box component="form" onSubmit={formik.handleSubmit} sx={{ mt: 2 }}>
      <TextField
        fullWidth
        id="content"
        name="content"
        label="Ваш комментарий"
        multiline
        rows={2}
        value={formik.values.content}
        onChange={formik.handleChange}
        error={formik.touched.content && Boolean(formik.errors.content)}
        helperText={formik.touched.content && formik.errors.content}
        sx={{ mb: 1 }}
      />
      <Button 
        type="submit" 
        variant="contained" 
        size="small"
        disabled={!formik.isValid || formik.isSubmitting}
      >
        Отправить
      </Button>
    </Box>
  );
};

export default CommentForm; 
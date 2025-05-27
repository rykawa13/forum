import React from 'react';
import { Alert, AlertTitle } from '@mui/material';

const Error = ({ message }) => {
  if (!message) return null;

  return (
    <Alert severity="error" sx={{ mb: 2 }}>
      <AlertTitle>Ошибка</AlertTitle>
      {message}
    </Alert>
  );
};

export default Error;
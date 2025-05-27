import React from 'react';
import { CircularProgress, Box, Typography } from '@mui/material';

const Loading = ({ message }) => {
  return (
    <Box sx={{ 
      display: 'flex', 
      flexDirection: 'column', 
      alignItems: 'center', 
      justifyContent: 'center', 
      p: 4 
    }}>
      <CircularProgress sx={{ mb: 2 }} />
      {message && <Typography variant="body1">{message}</Typography>}
    </Box>
  );
};

export default Loading;
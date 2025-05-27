import React from 'react';
import { Box, Typography, Container } from '@mui/material';

const Footer = () => {
  return (
    <Box
      component="footer"
      sx={{
        py: 3,
        px: 2,
        mt: 'auto',
        backgroundColor: (theme) => theme.palette.grey[100],
        borderTop: (theme) => `1px solid ${theme.palette.divider}`
      }}
    >
      <Container maxWidth="lg">
        <Typography variant="body2" color="text.secondary" align="center">
          © {new Date().getFullYear()} Форум. Все права защищены.
        </Typography>
      </Container>
    </Box>
  );
};

export default Footer;
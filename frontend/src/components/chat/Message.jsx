import React from 'react';
import { useSelector } from 'react-redux';
import { Box, Typography, Paper } from '@mui/material';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';

const Message = ({ message }) => {
  const { user } = useSelector(state => state.auth);
  const isCurrentUser = user?.id === message.userId;

  return (
    <Box
      sx={{
        display: 'flex',
        justifyContent: isCurrentUser ? 'flex-end' : 'flex-start',
        mb: 2,
      }}
    >
      <Paper
        elevation={3}
        sx={{
          p: 2,
          maxWidth: '70%',
          bgcolor: isCurrentUser ? 'primary.main' : 'background.paper',
          color: isCurrentUser ? 'primary.contrastText' : 'text.primary',
        }}
      >
        <Typography variant="subtitle2" sx={{ fontWeight: 'bold' }}>
          {message.username}
        </Typography>
        <Typography variant="body1">{message.text}</Typography>
        <Typography variant="caption" sx={{ display: 'block', textAlign: 'right' }}>
          {format(new Date(message.createdAt), 'HH:mm, dd MMMM', { locale: ru })}
        </Typography>
      </Paper>
    </Box>
  );
};

export default Message;
import React from 'react';
import { Box, Typography, Avatar } from '@mui/material';
import { format } from 'date-fns';
import { ru } from 'date-fns/locale';
import '../../styles/chat.css';

const Message = ({ message, isOwnMessage }) => {
  if (!message) return null;

  const {
    content = '',
    username = 'Anonymous',
    created_at
  } = message;

  const formattedTime = created_at ? format(new Date(created_at), 'HH:mm', { locale: ru }) : '';
  const firstLetter = username ? username.charAt(0).toUpperCase() : 'A';

  return (
    <Box 
      className={`message ${isOwnMessage ? 'own-message' : ''}`}
      sx={{
        display: 'flex',
        flexDirection: 'column',
        mb: 2,
        maxWidth: '70%',
        alignSelf: isOwnMessage ? 'flex-end' : 'flex-start',
        bgcolor: isOwnMessage ? 'primary.light' : 'background.paper',
        borderRadius: 2,
        p: 1
      }}
    >
      <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
        <Avatar 
          sx={{ 
            width: 24, 
            height: 24, 
            mr: 1,
            bgcolor: isOwnMessage ? 'primary.main' : 'secondary.main',
            fontSize: '0.875rem'
          }}
        >
          {firstLetter}
        </Avatar>
        <Typography 
          variant="subtitle2" 
          component="span"
          sx={{ 
            fontWeight: 'bold',
            color: isOwnMessage ? 'primary.main' : 'secondary.main'
          }}
        >
          {username}
        </Typography>
        {formattedTime && (
          <Typography 
            variant="caption" 
            component="span"
            sx={{ ml: 'auto', color: 'text.secondary' }}
          >
            {formattedTime}
          </Typography>
        )}
      </Box>
      <Typography variant="body1" sx={{ wordBreak: 'break-word' }}>
        {content}
      </Typography>
    </Box>
  );
};

export default Message;
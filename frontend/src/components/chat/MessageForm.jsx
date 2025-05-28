import React, { useState } from 'react';
import { TextField, Button, Box } from '@mui/material';
import SendIcon from '@mui/icons-material/Send';

const MessageForm = ({ onSendMessage, disabled }) => {
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message.trim() && !disabled) {
      onSendMessage(message);
      setMessage('');
    }
  };

  return (
    <Box 
      component="form" 
      onSubmit={handleSubmit} 
      sx={{ 
        display: 'flex', 
        p: 1,
        bgcolor: 'background.paper',
        borderTop: 1,
        borderColor: 'divider'
      }}
    >
      <TextField
        fullWidth
        variant="outlined"
        placeholder={disabled ? "Подключение к чату..." : "Напишите сообщение..."}
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        disabled={disabled}
        sx={{ mr: 1 }}
        size="small"
      />
      <Button
        type="submit"
        variant="contained"
        color="primary"
        disabled={!message.trim() || disabled}
        endIcon={<SendIcon />}
      >
        Отправить
      </Button>
    </Box>
  );
};

export default MessageForm;
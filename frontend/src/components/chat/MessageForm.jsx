import React, { useState } from 'react';
import { TextField, Button, Box } from '@mui/material';
import SendIcon from '@mui/icons-material/Send';

const MessageForm = ({ onSendMessage }) => {
  const [message, setMessage] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (message.trim()) {
      onSendMessage(message);
      setMessage('');
    }
  };

  return (
    <Box component="form" onSubmit={handleSubmit} sx={{ display: 'flex', p: 1 }}>
      <TextField
        fullWidth
        variant="outlined"
        placeholder="Напишите сообщение..."
        value={message}
        onChange={(e) => setMessage(e.target.value)}
        sx={{ mr: 1 }}
      />
      <Button
        type="submit"
        variant="contained"
        color="primary"
        disabled={!message.trim()}
        endIcon={<SendIcon />}
      >
        Отправить
      </Button>
    </Box>
  );
};

export default MessageForm;
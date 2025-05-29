import React, { useEffect, useCallback, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { connectSocket, disconnectSocket } from '../utils/socket';
import { fetchMessages, sendMessage, receiveMessage } from '../store/chatSlice';
import ChatWindow from '../components/chat/ChatWindow';
import { Box, Typography, Alert } from '@mui/material';

const ChatPage = () => {
  const dispatch = useDispatch();
  const { isAuthenticated } = useSelector(state => state.auth);
  const { messages } = useSelector(state => state.chat);
  const [connectionInfo, setConnectionInfo] = useState('');

  const handleNewMessage = useCallback((message) => {
    console.log('Handling new message in ChatPage:', message);
    if (message.type === 'connection_info') {
      setConnectionInfo(message.error);
    } else {
      dispatch(receiveMessage(message));
    }
  }, [dispatch]);

  useEffect(() => {
    console.log('ChatPage: Initializing chat connection');
    const token = localStorage.getItem('token');
    
    // Подключаемся к WebSocket с токеном или без него
    console.log('ChatPage: Connecting to WebSocket');
    connectSocket(token || '', handleNewMessage);

    return () => {
      console.log('ChatPage: Cleaning up chat connection');
      disconnectSocket();
    };
  }, [handleNewMessage]);

  const handleSendMessage = (messageText) => {
    if (!isAuthenticated) {
      setConnectionInfo('Для отправки сообщений необходима авторизация');
      return;
    }
    
    if (messageText.trim()) {
      console.log('ChatPage: Sending message:', messageText);
      dispatch(sendMessage(messageText));
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Чат форума
      </Typography>
      {connectionInfo && (
        <Alert severity="info" sx={{ mb: 2 }}>
          {connectionInfo}
        </Alert>
      )}
      <ChatWindow 
        messages={messages} 
        onSendMessage={handleSendMessage}
        isReadOnly={!isAuthenticated}
      />
    </Box>
  );
};

export default ChatPage;
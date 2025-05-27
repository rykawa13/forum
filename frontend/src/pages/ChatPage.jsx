import React, { useEffect } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { connectSocket, disconnectSocket } from '../utils/socket';
import { fetchMessages, sendMessage } from '../store/chatSlice';
import ChatWindow from '../components/chat/ChatWindow';
import { Box, Typography } from '@mui/material';

const ChatPage = () => {
  const dispatch = useDispatch();
  const { isAuthenticated, user } = useSelector(state => state.auth);
  const { messages, onlineUsers } = useSelector(state => state.chat);

  useEffect(() => {
    dispatch(fetchMessages());

    if (isAuthenticated) {
      connectSocket(user.token);
    }

    return () => {
      disconnectSocket();
    };
  }, [dispatch, isAuthenticated, user]);

  const handleSendMessage = (messageText) => {
    if (messageText.trim()) {
      dispatch(sendMessage({ text: messageText }));
    }
  };

  return (
    <Box sx={{ p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Чат форума
      </Typography>
      <Typography variant="subtitle1" gutterBottom>
        Онлайн: {onlineUsers} пользователей
      </Typography>
      <ChatWindow messages={messages} onSendMessage={handleSendMessage} />
    </Box>
  );
};

export default ChatPage;
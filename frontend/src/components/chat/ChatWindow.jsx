import React, { useEffect, useRef, useState } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchMessages, sendMessage, clearError } from '../../store/chatSlice';
import Message from './Message';
import MessageForm from './MessageForm';
import { Alert, Snackbar, Box, CircularProgress, Typography, IconButton } from '@mui/material';
import KeyboardArrowDownIcon from '@mui/icons-material/KeyboardArrowDown';
import '../../styles/chat.css';

const ChatWindow = () => {
  const dispatch = useDispatch();
  const { messages = [], loading, error, connected, connecting } = useSelector(state => state.chat);
  const { isAuthenticated, user } = useSelector(state => state.auth);
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);
  const [showScrollButton, setShowScrollButton] = useState(false);

  useEffect(() => {
    if (isAuthenticated) {
      console.log('ChatWindow: Fetching initial messages');
      dispatch(fetchMessages());
    }
  }, [dispatch, isAuthenticated]);

  useEffect(() => {
    const container = messagesContainerRef.current;
    if (!container) return;

    const handleScroll = () => {
      const { scrollTop, scrollHeight, clientHeight } = container;
      const isNearBottom = scrollHeight - scrollTop - clientHeight < 100;
      setShowScrollButton(!isNearBottom);
    };

    container.addEventListener('scroll', handleScroll);
    return () => container.removeEventListener('scroll', handleScroll);
  }, []);

  useEffect(() => {
    if (messages.length > 0) {
      const container = messagesContainerRef.current;
      if (!container) return;

      const { scrollTop, scrollHeight, clientHeight } = container;
      const isNearBottom = scrollHeight - scrollTop - clientHeight < 100;

      if (isNearBottom) {
        scrollToBottom();
      }
    }
  }, [messages]);

  const scrollToBottom = () => {
    if (messagesEndRef.current) {
      messagesEndRef.current.scrollIntoView({ behavior: 'smooth' });
    }
  };

  const handleScrollToBottom = () => {
    scrollToBottom();
  };

  const handleSendMessage = async (message) => {
    if (message.trim() && isAuthenticated) {
      if (!connected) {
        console.error('ChatWindow: Cannot send message - not connected');
        return;
      }
      console.log('ChatWindow: Sending message:', message);
      try {
        await dispatch(sendMessage(message)).unwrap();
        scrollToBottom();
      } catch (err) {
        console.error('ChatWindow: Error sending message:', err);
      }
    }
  };

  const handleCloseError = () => {
    dispatch(clearError());
  };

  if (!isAuthenticated) {
    return (
      <Box className="chat-window">
        <Box className="messages-container">
          <Typography className="info-message">
            Please log in to participate in the chat
          </Typography>
        </Box>
      </Box>
    );
  }

  if (loading) {
    return (
      <Box className="chat-window">
        <Box className="messages-container" sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
          <CircularProgress />
        </Box>
      </Box>
    );
  }

  return (
    <Box className="chat-window">
      <Snackbar 
        open={!!error} 
        autoHideDuration={6000} 
        onClose={handleCloseError}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert onClose={handleCloseError} severity="error" sx={{ width: '100%' }}>
          {error}
        </Alert>
      </Snackbar>

      {connecting && (
        <Box sx={{ p: 1, bgcolor: 'warning.light', textAlign: 'center' }}>
          <Typography>Connecting to chat...</Typography>
        </Box>
      )}

      {!connected && !connecting && (
        <Box sx={{ p: 1, bgcolor: 'error.light', textAlign: 'center' }}>
          <Typography>Chat disconnected. Trying to reconnect...</Typography>
        </Box>
      )}

      <Box className="messages-container" ref={messagesContainerRef}>
        {Array.isArray(messages) && messages.length > 0 ? (
          messages.map((msg) => (
            <Message 
              key={msg.id || Date.now()} 
              message={msg}
              isOwnMessage={msg.user_id === user?.id}
            />
          ))
        ) : (
          <Typography className="info-message">No messages yet</Typography>
        )}
        <div ref={messagesEndRef} />
      </Box>

      {showScrollButton && (
        <IconButton
          className="scroll-to-bottom"
          onClick={handleScrollToBottom}
          color="primary"
          sx={{
            position: 'absolute',
            bottom: 80,
            right: 20,
            backgroundColor: 'background.paper',
            boxShadow: 2,
            '&:hover': {
              backgroundColor: 'background.paper',
            },
          }}
        >
          <KeyboardArrowDownIcon />
        </IconButton>
      )}

      <MessageForm 
        onSendMessage={handleSendMessage} 
        disabled={!connected || connecting}
      />
    </Box>
  );
};

export default ChatWindow;
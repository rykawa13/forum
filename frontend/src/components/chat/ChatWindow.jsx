import React, { useEffect, useRef } from 'react';
import { useSelector, useDispatch } from 'react-redux';
import { fetchMessages, sendMessage, receiveMessage } from '../../store/chatSlice';
import Message from './Message';
import MessageForm from './MessageForm';
import '../../styles/chat.css';

const ChatWindow = () => {
  const dispatch = useDispatch();
  const { messages, loading, error } = useSelector(state => state.chat);
  const { isAuthenticated } = useSelector(state => state.auth);
  const messagesEndRef = useRef(null);

  useEffect(() => {
    dispatch(fetchMessages());
  }, [dispatch]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = (message) => {
    if (message.trim()) {
      dispatch(sendMessage(message));
    }
  };

  if (loading) return <div className="loading">Loading messages...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <div className="chat-window">
      <div className="messages-container">
        {messages.map((msg) => (
          <Message key={msg.id} message={msg} />
        ))}
        <div ref={messagesEndRef} />
      </div>
      {isAuthenticated && <MessageForm onSendMessage={handleSendMessage} />}
    </div>
  );
};

export default ChatWindow;
import React, { useState, useEffect } from 'react';
import { useWebSocket } from '../../hooks/useWebSocket';
import { fetchMessages } from '../../services/chatService';
import { useAuth } from '../../hooks/useAuth';

const Chat = () => {
  const [messages, setMessages] = useState([]);
  const [input, setInput] = useState('');
  const { sendMessage } = useWebSocket('ws://localhost/chat');
  const { isAuthenticated } = useAuth();

  useEffect(() => {
    const loadHistory = async () => {
      try {
        const history = await fetchMessages();
        setMessages(history);
      } catch (error) {
        console.error('Error loading chat history:', error);
      }
    };
    loadHistory();
  }, []);

  const handleSubmit = (e) => {
    e.preventDefault();
    if (!isAuthenticated) {
      alert('Please login to send messages');
      return;
    }
    if (input.trim()) {
      sendMessage({ 
        text: input, 
        timestamp: new Date().toISOString(),
        user: localStorage.getItem('username') || 'Anonymous'
      });
      setInput('');
    }
  };

  return (
    <div className="chat-container">
      <div className="messages">
        {messages.map((msg, i) => (
          <div key={i} className="message">
            <span className="message-time">
              {new Date(msg.timestamp).toLocaleTimeString()}
            </span>
            <strong>{msg.user}:</strong> {msg.text}
          </div>
        ))}
      </div>
      <form onSubmit={handleSubmit} className="message-form">
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
          placeholder="Type a message..."
          disabled={!isAuthenticated}
        />
        <button 
          type="submit" 
          disabled={!isAuthenticated}
          title={!isAuthenticated ? "Login to send messages" : ""}
        >
          Send
        </button>
      </form>
    </div>
  );
};

export default Chat;
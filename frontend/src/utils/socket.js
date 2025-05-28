import store from '../store/index';
import { setConnected, setConnecting } from '../store/chatSlice';

let socket = null;
let messageCallback;
let reconnectAttempts = 0;

export const connectSocket = (token, onMessage) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    console.log('Socket already connected');
    return;
  }

  store.dispatch(setConnecting(true));
  
  const wsUrl = `${process.env.REACT_APP_CHAT_WS_URL}?token=${token}`;
  console.log('Attempting to connect to WebSocket:', wsUrl);
  
  socket = new WebSocket(wsUrl);
  messageCallback = onMessage;
  
  socket.onopen = () => {
    console.log('WebSocket Connected');
    reconnectAttempts = 0;
    store.dispatch(setConnected(true));
    store.dispatch(setConnecting(false));
  };

  socket.onclose = () => {
    console.log('WebSocket Disconnected');
    store.dispatch(setConnected(false));
    
    if (reconnectAttempts < 5) {
      const delay = Math.min(1000 * Math.pow(2, reconnectAttempts), 30000);
      reconnectAttempts++;
      
      setTimeout(() => {
        connectSocket(token, onMessage);
      }, delay);
    }
  };

  socket.onerror = (error) => {
    console.error('WebSocket Error:', error);
  };

  socket.onmessage = (event) => {
    try {
      const message = JSON.parse(event.data);
      if (messageCallback) {
        messageCallback(message);
      }
    } catch (err) {
      console.error('Error processing message:', err);
    }
  };
};

export const disconnectSocket = () => {
  if (socket) {
    socket.close();
    socket = null;
    messageCallback = null;
    reconnectAttempts = 0;
    store.dispatch(setConnected(false));
    store.dispatch(setConnecting(false));
  }
};

export const sendMessage = (message) => {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    throw new Error('Socket not connected!');
  }
  
  console.log('Sending message:', message);
  socket.send(JSON.stringify(message));
};

export const isSocketConnected = () => {
  return socket && socket.readyState === WebSocket.OPEN;
};
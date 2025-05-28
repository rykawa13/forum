let socket = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 5;
const RECONNECT_DELAY = 2000;

export const connectSocket = (token) => {
  return new Promise((resolve, reject) => {
    try {
      if (socket && socket.readyState === WebSocket.OPEN) {
        console.log('WebSocket already connected');
        resolve(socket);
        return;
      }

      const wsUrl = `ws://localhost:8083/api/chat/ws?token=${token}`;
      console.log('Attempting to connect to WebSocket:', { url: wsUrl });

      socket = new WebSocket(wsUrl);

      socket.onopen = () => {
        console.log('WebSocket Connected', { readyState: socket.readyState });
        reconnectAttempts = 0;
        resolve(socket);
      };

      socket.onclose = (event) => {
        console.log('WebSocket Disconnected', { 
          code: event.code,
          reason: event.reason,
          wasClean: event.wasClean,
          readyState: socket?.readyState 
        });

        // Пытаемся переподключиться только если это не было чистое закрытие
        if (!event.wasClean && reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
          console.log(`Attempting to reconnect (${reconnectAttempts + 1}/${MAX_RECONNECT_ATTEMPTS})`);
          setTimeout(() => {
            reconnectAttempts++;
            connectSocket(token).catch(console.error);
          }, RECONNECT_DELAY);
        } else if (reconnectAttempts >= MAX_RECONNECT_ATTEMPTS) {
          console.log('Max reconnection attempts reached', { attempts: reconnectAttempts });
          reject(new Error('Max reconnection attempts reached'));
        }
      };

      socket.onerror = (error) => {
        console.log('WebSocket Error:', { 
          error,
          readyState: socket?.readyState,
          bufferedAmount: socket?.bufferedAmount 
        });
      };

      socket.onmessage = (event) => {
        try {
          console.log('Received message:', { data: event.data });
          const message = JSON.parse(event.data);
          
          if (message.type === 'auth_success') {
            console.log('Authentication successful');
          }
          
          // Добавляем обработчик для пинг-сообщений
          if (message.type === 'ping') {
            socket.send(JSON.stringify({ type: 'pong' }));
          }
        } catch (err) {
          console.error('Error processing message:', err);
        }
      };

    } catch (error) {
      console.error('Error creating WebSocket:', error);
      reject(error);
    }
  });
};

export const disconnectSocket = () => {
  if (socket) {
    console.log('Disconnecting WebSocket...', { readyState: socket.readyState });
    socket.close();
    socket = null;
    reconnectAttempts = 0;
  }
};

export const sendMessage = (message) => {
  if (socket && socket.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify(message));
  } else {
    console.error('Cannot send message - socket is not connected', { 
      socketExists: !!socket,
      readyState: socket?.readyState 
    });
  }
}; 
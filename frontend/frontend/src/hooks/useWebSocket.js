import { useEffect, useRef, useCallback } from 'react';

export const useWebSocket = (url) => {
  const ws = useRef(null);
  const [messages, setMessages] = useState([]);

  const sendMessage = useCallback((message) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(message));
    }
  }, []);

  useEffect(() => {
    ws.current = new WebSocket(url);
    
    ws.current.onmessage = (event) => {
      const message = JSON.parse(event.data);
      setMessages(prev => [...prev, message]);
    };

    return () => {
      if (ws.current?.readyState === WebSocket.OPEN) {
        ws.current.close();
      }
    };
  }, [url]);

  return { messages, sendMessage };
};
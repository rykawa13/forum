const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = function(app) {
  // Прокси для чат-сервиса
  const chatProxy = createProxyMiddleware({
    target: 'http://localhost:8083',
    changeOrigin: true,
    ws: true,
    pathRewrite: {
      '^/api/chat': '/api/chat',
      '^/ws': '/api/chat/ws'
    }
  });

  // Прокси для основного API
  const apiProxy = createProxyMiddleware({
    target: 'http://localhost:8081',
    changeOrigin: true,
  });

  app.use('/api/chat', chatProxy);
  app.use('/ws', chatProxy);
  app.use('/api', apiProxy);
};
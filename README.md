# Forum Project

Проект форума с микросервисной архитектурой.

## Сервисы

### Frontend
- React.js приложение
- Material-UI для компонентов интерфейса
- Redux для управления состоянием
- React Router для маршрутизации

### Backend

#### Auth Service (порт 8081)
- Аутентификация и авторизация пользователей
- Управление пользователями
- JWT токены

#### Forum Service (порт 8082)
- Управление постами и комментариями
- Взаимодействие с базой данных форума

#### Chat Service (порт 8083)
- WebSocket для real-time чата
- Хранение истории сообщений

## Установка и запуск

### Требования
- Node.js v16+
- Go 1.19+
- PostgreSQL 14+

### Frontend
```bash
cd frontend
npm install
npm start
```

### Backend

#### Auth Service
```bash
cd backend/auth-service
go mod download
go run cmd/main.go
```

#### Forum Service
```bash
cd backend/forum-service
go mod download
go run cmd/main.go
```

#### Chat Service
```bash
cd backend/chat-service
go mod download
go run cmd/main.go
```

## Конфигурация

Каждый сервис требует свой файл .env с необходимыми переменными окружения. Примеры конфигурации находятся в соответствующих директориях в файлах .env.example.

## API Documentation

Swagger документация доступна по следующим адресам:
- Auth Service: http://localhost:8081/swagger/index.html
- Forum Service: http://localhost:8082/swagger/index.html
- Chat Service: http://localhost:8083/swagger/index.html 
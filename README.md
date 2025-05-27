# Auth Service

Сервис аутентификации и авторизации с административной панелью.

## Технологии

### Backend
- Go
- Gin Web Framework
- PostgreSQL
- JWT для аутентификации
- Swagger для документации API

### Frontend
- React
- Axios для HTTP-запросов
- Material-UI для компонентов

## Функциональность

- Регистрация и авторизация пользователей
- JWT аутентификация
- Административная панель
- Управление пользователями
- Статистика форума
- CORS поддержка

## Установка и запуск

### Backend

```bash
cd backend/auth-service
go mod download
go run cmd/main.go
```

### Frontend

```bash
cd frontend
npm install
npm start
```

Сервер запускается на `http://localhost:8081`
Клиент запускается на `http://localhost:3000`

## API Документация

Swagger документация доступна по адресу: `http://localhost:8081/swagger/index.html` 
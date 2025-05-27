#!/bin/sh

# Ждем доступности PostgreSQL
echo "Waiting for PostgreSQL..."
while ! pg_isready -h postgres -p 5432 -U forum_user
do
    sleep 2
done

# Запускаем миграции
echo "Running migrations..."
migrate -path /app/migrations -database "$DB_URL" up

# Запускаем приложение
echo "Starting application..."
/app/main 
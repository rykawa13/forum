#!/bin/bash

# Установка swag, если он еще не установлен
if ! command -v swag &> /dev/null; then
    echo "Installing swag..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Генерация документации
echo "Generating Swagger documentation..."
swag init -g cmd/app/main.go -o docs

echo "Done!" 
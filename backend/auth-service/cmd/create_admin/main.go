package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type UserCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"is_admin"`
}

func main() {
	// Создаем данные для запроса
	admin := UserCreate{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "admin",
		IsAdmin:  true,
	}

	// Преобразуем структуру в JSON
	jsonData, err := json.Marshal(admin)
	if err != nil {
		fmt.Printf("Ошибка при создании JSON: %v\n", err)
		return
	}

	// Отправляем POST запрос
	resp, err := http.Post(
		"http://localhost:8081/auth/sign-up",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		fmt.Printf("Ошибка при отправке запроса: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusOK {
		fmt.Println("Админ успешно создан!")
		fmt.Println("Теперь вы можете войти с данными:")
		fmt.Println("Email: admin@example.com")
		fmt.Println("Password: admin")
	} else {
		// Читаем тело ответа для получения деталей ошибки
		var response map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Printf("Ошибка при создании админа. Статус: %d\n", resp.StatusCode)
		} else {
			fmt.Printf("Ошибка при создании админа: %v\n", response["error"])
		}
	}
}

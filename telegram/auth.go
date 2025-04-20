package telegram

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

type InitData struct {
	QueryID  string `json:"query_id"`
	User     User   `json:"user"`
	AuthDate int64  `json:"auth_date"`
	Hash     string `json:"hash"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func ValidateInitData(initData string) (*InitData, error) {
	// Парсим данные инициализации
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, fmt.Errorf("error parsing init data: %v", err)
	}

	// Получаем хеш
	hash := values.Get("hash")
	if hash == "" {
		return nil, fmt.Errorf("hash not found in init data")
	}

	// Удаляем хеш из данных для проверки
	values.Del("hash")

	// Сортируем ключи
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Создаем строку для проверки
	var dataCheckString strings.Builder
	for _, k := range keys {
		if dataCheckString.Len() > 0 {
			dataCheckString.WriteString("\n")
		}
		dataCheckString.WriteString(k)
		dataCheckString.WriteString("=")
		dataCheckString.WriteString(values.Get(k))
	}

	// Получаем токен бота
	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN not set")
	}

	// Создаем секретный ключ
	secretKey := sha256.Sum256([]byte("WebAppData"))

	// Создаем HMAC
	mac := hmac.New(sha256.New, secretKey[:])
	mac.Write([]byte(dataCheckString.String()))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	// Проверяем хеш
	if hash != expectedHash {
		return nil, fmt.Errorf("invalid hash")
	}

	// Проверяем время авторизации
	authDate := values.Get("auth_date")
	if authDate == "" {
		return nil, fmt.Errorf("auth_date not found")
	}

	authTime, err := time.Parse("2006-01-02 15:04:05", authDate)
	if err != nil {
		return nil, fmt.Errorf("invalid auth_date format: %v", err)
	}

	// Проверяем, что авторизация не старше 24 часов
	if time.Since(authTime) > 24*time.Hour {
		return nil, fmt.Errorf("auth_date expired")
	}

	// Возвращаем данные пользователя
	return &InitData{
		QueryID:  values.Get("query_id"),
		User:     User{}, // Заполнить из values
		AuthDate: authTime.Unix(),
		Hash:     hash,
	}, nil
}

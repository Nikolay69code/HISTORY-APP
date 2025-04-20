package db

import "time"

// User представляет пользователя в системе
type User struct {
	ID         int64     `json:"id"`
	TelegramID int64     `json:"telegram_id"`
	Username   string    `json:"username"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
}

type Topic struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	OrderNum    int    `json:"order_num"`
}

// Task представляет задание
type Task struct {
	ID          int      `json:"id"`
	TopicID     int      `json:"topic_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Difficulty  int      `json:"difficulty"`
	Options     []Option `json:"options"`
}

// Option представляет вариант ответа
type Option struct {
	ID        int    `json:"id"`
	Text      string `json:"text"`
	IsCorrect bool   `json:"is_correct"`
}

// TheoryMaterial представляет теоретический материал
type TheoryMaterial struct {
	ID       int    `json:"id"`
	TopicID  int    `json:"topic_id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	OrderNum int    `json:"order_num"`
}

type UserProgress struct {
	ID            int       `json:"id"`
	UserID        int64     `json:"user_id"`
	TaskID        int       `json:"task_id"`
	IsCorrect     bool      `json:"is_correct"`
	AttemptCount  int       `json:"attempt_count"`
	LastAttemptAt time.Time `json:"last_attempt_at"`
}

// Statistics представляет статистику пользователя
type Statistics struct {
	TopicID         int `json:"topic_id"`
	TotalAttempts   int `json:"total_attempts"`
	CorrectAttempts int `json:"correct_attempts"`
}

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error connecting to the database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error pinging the database: %v", err)
	}

	log.Println("Successfully connected to database")
	return nil
}

// GetUserByTelegramID получает пользователя по его Telegram ID
func GetUserByTelegramID(telegramID int64) (*User, error) {
	user := &User{}
	err := DB.QueryRow(`
		SELECT id, telegram_id, username, first_name, last_name, created_at 
		FROM users 
		WHERE telegram_id = $1
	`, telegramID).Scan(&user.ID, &user.TelegramID, &user.Username, &user.FirstName, &user.LastName, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// CreateUser создает нового пользователя
func CreateUser(user *User) error {
	_, err := DB.Exec(`
		INSERT INTO users (telegram_id, username, first_name, last_name)
		VALUES ($1, $2, $3, $4)
	`, user.TelegramID, user.Username, user.FirstName, user.LastName)
	return err
}

// GetTasksByTopic получает задания по теме
func GetTasksByTopic(topicID int) ([]Task, error) {
	rows, err := DB.Query(`
		SELECT t.id, t.title, t.description, t.difficulty,
			   o.id, o.text, o.is_correct
		FROM tasks t
		LEFT JOIN options o ON o.task_id = t.id
		WHERE t.topic_id = $1
		ORDER BY t.id, o.id
	`, topicID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]Task, 0)
	currentTask := Task{}

	for rows.Next() {
		var optionID int
		var optionText string
		var isCorrect bool

		err := rows.Scan(
			&currentTask.ID,
			&currentTask.Title,
			&currentTask.Description,
			&currentTask.Difficulty,
			&optionID,
			&optionText,
			&isCorrect,
		)
		if err != nil {
			return nil, err
		}

		option := Option{
			ID:        optionID,
			Text:      optionText,
			IsCorrect: isCorrect,
		}
		currentTask.Options = append(currentTask.Options, option)

		if len(tasks) == 0 || tasks[len(tasks)-1].ID != currentTask.ID {
			tasks = append(tasks, currentTask)
			currentTask = Task{}
		}
	}

	return tasks, nil
}

// SaveProgress сохраняет прогресс пользователя
func SaveProgress(userID int64, taskID int, isCorrect bool) error {
	_, err := DB.Exec(`
		INSERT INTO user_progress (user_id, task_id, is_correct)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, task_id)
		DO UPDATE SET 
			attempt_count = user_progress.attempt_count + 1,
			is_correct = $3,
			last_attempt_at = CURRENT_TIMESTAMP
	`, userID, taskID, isCorrect)
	return err
}

// GetUserStatistics получает статистику пользователя
func GetUserStatistics(userID int64) ([]Statistics, error) {
	rows, err := DB.Query(`
		SELECT topic_id, total_attempts, correct_attempts
		FROM statistics
		WHERE user_id = $1
		ORDER BY topic_id
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stats []Statistics
	for rows.Next() {
		var stat Statistics
		err := rows.Scan(&stat.TopicID, &stat.TotalAttempts, &stat.CorrectAttempts)
		if err != nil {
			return nil, err
		}
		stats = append(stats, stat)
	}
	return stats, nil
}

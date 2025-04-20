package db

import (
	"database/sql"
	"fmt"
)

// GetNextTaskForUser получает следующее задание для пользователя
func GetNextTaskForUser(userID int64) (*Task, error) {
	// Получаем задание, которое пользователь еще не решал или решал неправильно
	query := `
		SELECT t.id, t.topic_id, t.title, t.description, t.difficulty
		FROM tasks t
		LEFT JOIN user_progress up ON t.id = up.task_id AND up.user_id = $1
		WHERE up.id IS NULL OR (up.is_correct = false)
		ORDER BY t.difficulty, RANDOM()
		LIMIT 1
	`

	task := &Task{}
	err := DB.QueryRow(query, userID).Scan(
		&task.ID,
		&task.TopicID,
		&task.Title,
		&task.Description,
		&task.Difficulty,
	)

	if err == sql.ErrNoRows {
		// Если все задания решены, возвращаем случайное
		query = `
			SELECT id, topic_id, title, description, difficulty
			FROM tasks
			ORDER BY RANDOM()
			LIMIT 1
		`
		err = DB.QueryRow(query).Scan(
			&task.ID,
			&task.TopicID,
			&task.Title,
			&task.Description,
			&task.Difficulty,
		)
	}

	if err != nil {
		return nil, fmt.Errorf("error getting next task: %v", err)
	}

	// Получаем варианты ответов для задания
	query = `
		SELECT id, text, is_correct
		FROM options
		WHERE task_id = $1
	`
	rows, err := DB.Query(query, task.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting task options: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var option Option
		err := rows.Scan(&option.ID, &option.Text, &option.IsCorrect)
		if err != nil {
			return nil, fmt.Errorf("error scanning option: %v", err)
		}
		task.Options = append(task.Options, option)
	}

	return task, nil
}

// GetTheoryMaterialsByTopic получает теоретические материалы по теме
func GetTheoryMaterialsByTopic(topicID string) ([]TheoryMaterial, error) {
	query := `
		SELECT id, topic_id, title, content, order_num
		FROM theory_materials
		WHERE topic_id = $1
		ORDER BY order_num
	`

	rows, err := DB.Query(query, topicID)
	if err != nil {
		return nil, fmt.Errorf("error getting theory materials: %v", err)
	}
	defer rows.Close()

	var materials []TheoryMaterial
	for rows.Next() {
		var material TheoryMaterial
		err := rows.Scan(
			&material.ID,
			&material.TopicID,
			&material.Title,
			&material.Content,
			&material.OrderNum,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning theory material: %v", err)
		}
		materials = append(materials, material)
	}

	return materials, nil
}

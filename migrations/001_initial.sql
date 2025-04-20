-- Создание таблицы пользователей
CREATE TABLE users (
    id BIGINT PRIMARY KEY,
    telegram_id BIGINT UNIQUE NOT NULL,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы с темами
CREATE TABLE topics (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    order_num INT NOT NULL
);

-- Создание таблицы с заданиями
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    topic_id INT REFERENCES topics(id),
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    difficulty INT CHECK (difficulty BETWEEN 1 AND 5),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы с вариантами ответов
CREATE TABLE options (
    id SERIAL PRIMARY KEY,
    task_id INT REFERENCES tasks(id),
    text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL
);

-- Создание таблицы с теоретическими материалами
CREATE TABLE theory_materials (
    id SERIAL PRIMARY KEY,
    topic_id INT REFERENCES topics(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    order_num INT NOT NULL
);

-- Создание таблицы с прогрессом пользователей
CREATE TABLE user_progress (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    task_id INT REFERENCES tasks(id),
    is_correct BOOLEAN NOT NULL,
    attempt_count INT DEFAULT 1,
    last_attempt_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, task_id)
);

-- Создание таблицы с статистикой
CREATE TABLE statistics (
    id SERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    topic_id INT REFERENCES topics(id),
    total_attempts INT DEFAULT 0,
    correct_attempts INT DEFAULT 0,
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, topic_id)
); 
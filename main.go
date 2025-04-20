package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"history-ege-app/db"
	"history-ege-app/telegram"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Инициализация базы данных
	if err := db.Init(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := gin.Default()

	// Middleware для проверки авторизации
	authMiddleware := func(c *gin.Context) {
		initData := c.Query("initData")
		if initData == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No init data provided"})
			c.Abort()
			return
		}

		data, err := telegram.ValidateInitData(initData)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		// Получаем или создаем пользователя
		user, err := db.GetUserByTelegramID(data.User.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}

		if user == nil {
			user = &db.User{
				TelegramID: data.User.ID,
				Username:   data.User.Username,
				FirstName:  data.User.FirstName,
				LastName:   data.User.LastName,
			}
			if err := db.CreateUser(user); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				c.Abort()
				return
			}
		}

		c.Set("user", user)
		c.Next()
	}

	// Статические файлы
	r.Static("/static", "./frontend")
	r.LoadHTMLGlob("frontend/*.html")

	// Главная страница
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	// API маршруты
	api := r.Group("/api")
	api.Use(authMiddleware)
	{
		api.GET("/tasks", getTasks)
		api.GET("/next-task", getNextTask)
		api.POST("/save-progress", saveProgress)
		api.GET("/statistics", getStatistics)
		api.GET("/theory", getTheoryMaterials)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func getTasks(c *gin.Context) {
	topicIDStr := c.Query("topic_id")
	if topicIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic ID is required"})
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID format"})
		return
	}

	tasks, err := db.GetTasksByTopic(topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func getNextTask(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	task, err := db.GetNextTaskForUser(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get next task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func saveProgress(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	var progress struct {
		TaskID    int  `json:"task_id"`
		IsCorrect bool `json:"is_correct"`
	}

	if err := c.BindJSON(&progress); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := db.SaveProgress(user.ID, progress.TaskID, progress.IsCorrect); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func getStatistics(c *gin.Context) {
	user := c.MustGet("user").(*db.User)
	stats, err := db.GetUserStatistics(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get statistics"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func getTheoryMaterials(c *gin.Context) {
	topicIDStr := c.Query("topic_id")
	if topicIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Topic ID is required"})
		return
	}

	topicID, err := strconv.Atoi(topicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid topic ID format"})
		return
	}

	materials, err := db.GetTheoryMaterialsByTopic(topicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get theory materials"})
		return
	}

	c.JSON(http.StatusOK, materials)
}

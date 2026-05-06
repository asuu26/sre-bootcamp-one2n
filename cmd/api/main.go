package main

import (
	"os"

	"github.com/75asu/sre-bootcamp-one2n/internal/api/handlers"
	"github.com/75asu/sre-bootcamp-one2n/internal/db"
	"github.com/75asu/sre-bootcamp-one2n/internal/logger"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	log := logger.New()
	defer log.Sync()

	if err := godotenv.Load(); err != nil {
		log.Warn("no .env file found, reading from environment")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("database unreachable", zap.Error(err))
	}
	log.Info("database connected")

	gin.SetMode(os.Getenv("GIN_MODE"))

	r := gin.New()
	r.Use(logger.GinLogger(log))
	r.Use(gin.Recovery())
	r.SetTrustedProxies([]string{"127.0.0.1"})

	r.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	student := handlers.NewStudentHandler(conn, log)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/students", student.Create)
		v1.GET("/students", student.GetAll)
		v1.GET("/students/:id", student.GetByID)
		v1.PUT("/students/:id", student.Update)
		v1.DELETE("/students/:id", student.Delete)
	}

	log.Info("starting server", zap.String("port", port))
	if err := r.Run(":" + port); err != nil {
		log.Fatal("server failed", zap.Error(err))
	}
}

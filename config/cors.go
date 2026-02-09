package config

import (
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	origins := []string{"http://localhost:8000"}

	// CORS_ORIGINS can be a comma-separated list, e.g. "https://your-app.onrender.com,http://localhost:3000"
	if extra := os.Getenv("CORS_ORIGINS"); extra != "" {
		origins = strings.Split(extra, ",")
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}

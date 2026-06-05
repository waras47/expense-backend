package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func registerMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))
}

func registerRoutes(r *gin.Engine, h *handlers) {
	api := r.Group("/api")
	h.category.RegisterRoutes(api.Group("/categories"))
	// TODO: Add the required routers
}

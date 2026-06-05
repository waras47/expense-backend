package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// registerMiddleware memasang middleware global (CORS, dll).
func registerMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))
}

// registerRoutes mendaftarkan semua route API.
// Tiap handler mendaftarkan route-nya sendiri.
func registerRoutes(r *gin.Engine, h *handlers) {
	api := r.Group("/api")
	h.category.RegisterRoutes(api.Group("/categories"))
	// tambah modul baru di sini, contoh:
	// h.expense.RegisterRoutes(api.Group("/expenses"))
}

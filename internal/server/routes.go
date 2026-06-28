package server

import (
	"expense-backend/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func registerMiddleware(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"},
	}))
}

func registerRoutes(r *gin.Engine, h *handlers) {
	/* SWAGGER */
	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	/* APP ROUTES */
	api := r.Group("/api")
	h.category.RegisterRoutes(api.Group("/categories"))
	h.income.RegisterRoutes(api.Group("/incomes"))
	h.expense.RegisterRoutes(api.Group("/expenses"))
	// TODO: Add the required routers
}

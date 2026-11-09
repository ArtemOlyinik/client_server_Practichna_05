package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/config"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/database"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/handler"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/middleware"

	_ "github.com/ArtemOlyinik/client_server_Practichna_05/docs"
)

// @title Incident Tracker API
// @version 1.0
// @description REST API для обліку технічних інцидентів та аварій.
// @description Підтримує реєстрацію інженерів, JWT-авторизацію та повний CRUD інцидентів.
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введіть токен у форматі: Bearer <JWT>
func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("[GORM] Помилка підключення до PostgreSQL: %v", err)
	}
	log.Printf("[GORM] Connected to PostgreSQL successfully (host=%s:%s, db=%s)", cfg.DBHost, cfg.DBPort, cfg.DBName)

	if err := database.Migrate(db); err != nil {
		log.Fatalf("[GORM] Помилка міграції: %v", err)
	}
	log.Println("[GORM] Auto-migration completed: 'users', 'incidents' tables are up to date.")

	authHandler := handler.NewAuthHandler(db, cfg.JWTSecret)
	incidentHandler := handler.NewIncidentHandler(db)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		incidents := api.Group("/incidents")
		{
			incidents.GET("", incidentHandler.List)
			incidents.GET("/:id", incidentHandler.Get)

			protected := incidents.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protected.POST("", incidentHandler.Create)
				protected.PATCH("/:id/status", incidentHandler.UpdateStatus)
				protected.DELETE("/:id", incidentHandler.Delete)
			}
		}
	}

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	log.Printf("[СЕРВЕР] Сервіс успішно запущено на http://localhost:%s", cfg.Port)
	log.Printf("[СЕРВЕР] Swagger UI доступний за адресою: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("[СЕРВЕР] Помилка запуску Gin: %v", err)
	}
}

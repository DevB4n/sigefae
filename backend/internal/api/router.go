package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sigefae/internal/db"
	"sigefae/internal/auth"
)

func New(database *gorm.DB) *gin.Engine {

	router := gin.Default()

	// ==========================
	// Health Check
	// ==========================

	router.GET("/ping", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})

	})

	// ==========================
	// Auth
	// ==========================

	authService := auth.New(database)
	authHandler := auth.NewHandler(authService)

	api := router.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
	}

	// ==========================
	// Rutas protegidas
	// ==========================

	protected := router.Group("/api")
	protected.Use(auth.Middleware(database))

	{

		protected.GET("/me", func(c *gin.Context) {

			user := c.MustGet("user").(db.Usuario)

			c.JSON(http.StatusOK, gin.H{
				"id":     user.ID,
				"nombre": user.Nombre,
				"email":  user.Email,
				"cargo":  user.Cargo,
				"rol":    user.Rol.Nombre,
			})

		})

		// ==========================
		// Solo Superadministrador
		// ==========================

		admin := protected.Group("")
		admin.Use(auth.RequireRole("Superadministrador"))

		{
			// admin.POST("/usuarios", userHandler.Create)
		}
	}

	return router
}
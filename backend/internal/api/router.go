package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sigefae/internal/user"
	"sigefae/internal/db"
	"sigefae/internal/auth"
	"sigefae/internal/role"
	"sigefae/internal/tipo_pago"
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

	userService := user.New(database)
	userHandler := user.NewHandler(userService)
	roleService := role.New(database)
	roleHandler := role.NewHandler(roleService)

	tipoPagoService := tipo_pago.New(database)
	tipoPagoHandler := tipo_pago.NewHandler(tipoPagoService)

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
			

			// ==========================
			// Usuarios
			// ==========================
			
			admin.POST("/usuarios", userHandler.Create)
			admin.GET("/usuarios", userHandler.List)
			admin.GET("/usuarios/:id", userHandler.GetByID)
			admin.PUT("/usuarios/:id", userHandler.Update)
			admin.PATCH("/usuarios/:id/activo", userHandler.UpdateStatus)
			admin.PATCH("/usuarios/:id/password", userHandler.UpdatePassword)
			
			// ==========================
			// Roles
			// ==========================
			
			admin.POST("/roles", roleHandler.Create)
			admin.GET("/roles", roleHandler.List)
			admin.PATCH("/roles/:id/activo", roleHandler.UpdateStatus)


			// ==========================
			// Tipos de Pago
			// ==========================
			admin.POST("/tipos-pago", tipoPagoHandler.Create)
			admin.GET("/tipos-pago", tipoPagoHandler.List)
			admin.PATCH("/tipos-pago/:id/activo", tipoPagoHandler.UpdateStatus)
		}
	}

	return router
}
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"sigefae/internal/area"
	"sigefae/internal/auth"
	"sigefae/internal/db"
	"sigefae/internal/metodo_pago"
	"sigefae/internal/moneda"
	"sigefae/internal/pais"
	"sigefae/internal/paso_ruta"
	"sigefae/internal/role"
	"sigefae/internal/ruta"
	"sigefae/internal/tipo_factura"
	"sigefae/internal/tipo_pago"
	"sigefae/internal/user"
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

	metodoPagoService := metodo_pago.New(database)
	metodoPagoHandler := metodo_pago.NewHandler(metodoPagoService)

	areaService := area.New(database)
	areaHandler := area.NewHandler(areaService)

	tipoFacturaService := tipo_factura.New(database)
	tipoFacturaHandler := tipo_factura.NewHandler(tipoFacturaService)

	rutaService := ruta.New(database)
	rutaHandler := ruta.NewHandler(rutaService)

	pasoRutaService := paso_ruta.New(database)
	pasoRutaHandler := paso_ruta.NewHandler(pasoRutaService)

	monedaService := moneda.New(database)
	monedaHandler := moneda.NewHandler(monedaService)

	paisService := pais.New(database)
	paisHandler := pais.NewHandler(paisService)

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
			admin.PATCH("/roles/:id", roleHandler.Update)
			admin.PATCH("/roles/:id/activo", roleHandler.UpdateStatus)

			// ==========================/pasos-ruta/:id/activo
			// Tipos de Pago
			// ==========================
			admin.POST("/tipos-pago", tipoPagoHandler.Create)
			admin.GET("/tipos-pago", tipoPagoHandler.List)
			admin.PATCH("/tipos-pago/:id", tipoPagoHandler.Update)
			admin.PATCH("/tipos-pago/:id/activo", tipoPagoHandler.UpdateStatus)

			// ==========================
			// Métodos de Pago
			// ==========================

			admin.POST("/metodos-pago", metodoPagoHandler.Create)
			admin.GET("/metodos-pago", metodoPagoHandler.List)
			admin.PATCH("/metodos-pago/:id", metodoPagoHandler.Update)
			admin.PATCH("/metodos-pago/:id/activo", metodoPagoHandler.UpdateStatus)

			// ==========================
			// Áreas
			// ==========================

			admin.POST("/areas", areaHandler.Create)
			admin.GET("/areas", areaHandler.List)
			admin.PATCH("/areas/:id/activo", areaHandler.UpdateStatus)
			admin.PATCH("/areas/:id", areaHandler.Update)

			// ==========================
			// Tipos de Factura
			// ==========================

			admin.POST("/tipos-factura", tipoFacturaHandler.Create)
			admin.GET("/tipos-factura", tipoFacturaHandler.List)
			admin.PUT("/tipos-factura/:id", tipoFacturaHandler.Update)
			admin.PATCH("/tipos-factura/:id/activo", tipoFacturaHandler.UpdateStatus)

			// ==========================
			// Rutas
			// ==========================

			admin.POST("/rutas", rutaHandler.Create)
			admin.GET("/rutas", rutaHandler.List)
			admin.PUT("/rutas/:id", rutaHandler.Update)
			admin.PATCH("/rutas/:id/activo", rutaHandler.UpdateStatus)

			// ==========================
			// Pasos de Ruta
			// ==========================

			admin.POST("/pasos-ruta", pasoRutaHandler.Create)
			admin.GET("/pasos-ruta", pasoRutaHandler.List)
			admin.PUT("/pasos-ruta/:id", pasoRutaHandler.Update)
			admin.PATCH("/pasos-ruta/:id/activo", pasoRutaHandler.UpdateStatus)

			// ==========================
			// Moneda
			// ==========================

			admin.POST("/moneda", monedaHandler.Create)
			admin.GET("/moneda", monedaHandler.List)
			admin.PUT("/moneda/:id", monedaHandler.Update)
			admin.PATCH("/moneda/:id/activo", monedaHandler.UpdateStatus)

			//===========================
			//pais
			//===========================

			admin.POST("/pais", paisHandler.Create)
			admin.GET("/pais", paisHandler.List)
			admin.PUT("/pais/:id", paisHandler.Update)
			admin.PATCH("/pais/:id/activo", paisHandler.UpdateStatus)
		}
	}

	return router
}

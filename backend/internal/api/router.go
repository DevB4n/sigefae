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
	"sigefae/internal/departamento"
	"sigefae/internal/municipio"
	"sigefae/internal/direccion"
	"sigefae/internal/actividad_economica"
	"sigefae/internal/contacto"
	"sigefae/internal/tipo_documento"
	"sigefae/internal/categoria_proveedor"
	"sigefae/internal/tipo_persona"
	"sigefae/internal/proveedor"
	"sigefae/internal/responsabilidad_fiscal"
	"sigefae/internal/receptor"
	"sigefae/internal/estado_correo"
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

	departamentoService := departamento.New(database)
	departamentoHandler := departamento.NewHandler(departamentoService)

	municipioService := municipio.New(database)
	municipioHandler := municipio.NewHandler(municipioService)

	direccionService := direccion.New(database)
	direccionHandler := direccion.NewHandler(direccionService)

	actividadEconomicaService := actividad_economica.New(database)
	actividadEconomicaHandler := actividad_economica.NewHandler(actividadEconomicaService)

	tipoDocumentoService := tipo_documento.New(database)
	tipoDocumentoHandler := tipo_documento.NewHandler(tipoDocumentoService)
	
	categoriaProveedorService := categoria_proveedor.New(database)
	categoriaProveedorHandler := categoria_proveedor.NewHandler(categoriaProveedorService)

	tipoPersonaService := tipo_persona.New(database)
	tipoPersonaHandler := tipo_persona.NewHandler(tipoPersonaService)

	proveedorService := proveedor.New(database)
	proveedorHandler := proveedor.NewHandler(proveedorService)

	contactoService := contacto.New(database)
	contactoHandler := contacto.NewHandler(contactoService)

	responsabilidadFiscalService := responsabilidad_fiscal.New(database)
	responsabilidadFiscalHandler := responsabilidad_fiscal.NewHandler(responsabilidadFiscalService)

	receptorService := receptor.New(database)
	receptorHandler := receptor.NewHandler(receptorService)

	estadoCorreoService := estado_correo.New(database)
	estadoCorreoHandler := estado_correo.NewHandler(estadoCorreoService)
	
	api := router.Group("/api")
	{
		api.POST("/auth/login", authHandler.Login)
	}

	// =========================
	// Rutas protegidas
	// =========================

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

		// =========================
		// Solo Superadministrador
		// =========================

		admin := protected.Group("")
		admin.Use(auth.RequireRole("Superadministrador"))

		{

			// =========================
			// Usuarios
			// =========================

			admin.POST("/usuarios", userHandler.Create)
			admin.GET("/usuarios", userHandler.List)
			admin.GET("/usuarios/:id", userHandler.GetByID)
			admin.PUT("/usuarios/:id", userHandler.Update)
			admin.PATCH("/usuarios/:id/activo", userHandler.UpdateStatus)
			admin.PATCH("/usuarios/:id/password", userHandler.UpdatePassword)

			// =========================
			// Roles
			// =========================

			admin.POST("/roles", roleHandler.Create)
			admin.GET("/roles", roleHandler.List)
			admin.PATCH("/roles/:id", roleHandler.Update)
			admin.PATCH("/roles/:id/activo", roleHandler.UpdateStatus)

			// =========================/pasos-ruta/:id/activo
			// Tipos de Pago
			// =========================
			admin.POST("/tipos-pago", tipoPagoHandler.Create)
			admin.GET("/tipos-pago", tipoPagoHandler.List)
			admin.PATCH("/tipos-pago/:id", tipoPagoHandler.Update)
			admin.PATCH("/tipos-pago/:id/activo", tipoPagoHandler.UpdateStatus)

			// =========================
			// Métodos de Pago
			// =========================

			admin.POST("/metodos-pago", metodoPagoHandler.Create)
			admin.GET("/metodos-pago", metodoPagoHandler.List)
			admin.PATCH("/metodos-pago/:id", metodoPagoHandler.Update)
			admin.PATCH("/metodos-pago/:id/activo", metodoPagoHandler.UpdateStatus)

			// =========================
			// Áreas
			// =========================

			admin.POST("/areas", areaHandler.Create)
			admin.GET("/areas", areaHandler.List)
			admin.PATCH("/areas/:id/activo", areaHandler.UpdateStatus)
			admin.PATCH("/areas/:id", areaHandler.Update)

			// =========================
			// Tipos de Factura
			// =========================

			admin.POST("/tipos-factura", tipoFacturaHandler.Create)
			admin.GET("/tipos-factura", tipoFacturaHandler.List)
			admin.PUT("/tipos-factura/:id", tipoFacturaHandler.Update)
			admin.PATCH("/tipos-factura/:id/activo", tipoFacturaHandler.UpdateStatus)

			// =========================
			// Rutas
			// =========================

			admin.POST("/rutas", rutaHandler.Create)
			admin.GET("/rutas", rutaHandler.List)
			admin.PUT("/rutas/:id", rutaHandler.Update)
			admin.PATCH("/rutas/:id/activo", rutaHandler.UpdateStatus)

			// =========================
			// Pasos de Ruta
			// =========================

			admin.POST("/pasos-ruta", pasoRutaHandler.Create)
			admin.GET("/pasos-ruta", pasoRutaHandler.List)
			admin.PUT("/pasos-ruta/:id", pasoRutaHandler.Update)
			admin.PATCH("/pasos-ruta/:id/activo", pasoRutaHandler.UpdateStatus)

			// =========================
			// Moneda
			// =========================

			admin.POST("/moneda", monedaHandler.Create)
			admin.GET("/moneda", monedaHandler.List)
			admin.PUT("/moneda/:id", monedaHandler.Update)
			admin.PATCH("/moneda/:id/activo", monedaHandler.UpdateStatus)

			// =========================
			// Pais
			// =========================

			admin.POST("/pais", paisHandler.Create)
			admin.GET("/pais", paisHandler.List)
			admin.PUT("/pais/:id", paisHandler.Update)
			admin.PATCH("/pais/:id/activo", paisHandler.UpdateStatus)

			// =========================
			// Departamento
			// =========================

			admin.POST("/departamento", departamentoHandler.Create)
			admin.GET("/departamento", departamentoHandler.List)
			admin.PUT("/departamento/:id", departamentoHandler.Update)
			admin.PATCH("/departamento/:id/activo", departamentoHandler.UpdateStatus)

			// =========================
			//	Municipio
			// =========================

			admin.POST("/municipio", municipioHandler.Create)
			admin.GET("/municipio", municipioHandler.List)
			admin.PUT("/municipio/:id", municipioHandler.Update)
			admin.PATCH("/municipio/:id/activo", municipioHandler.UpdateStatus)

			// =========================
			//	Direccion
			// =========================

			admin.POST("/direccion", direccionHandler.Create)
			admin.GET("/direccion", direccionHandler.List)
			admin.PUT("/direccion/:id", direccionHandler.Update)
			admin.PATCH("/direccion/:id/activo", direccionHandler.UpdateStatus)

			// =========================
			//	Actividad_Economica
			// =========================


			admin.POST("/actividad-economica", actividadEconomicaHandler.Create)
			admin.GET("/actividad-economica", actividadEconomicaHandler.List)
			admin.PUT("/actividad-economica/:id", actividadEconomicaHandler.Update)
			admin.PATCH("/actividad-economica/:id/activo", actividadEconomicaHandler.UpdateStatus)

			// =========================
			//	Tipo_documento
			// =========================

			admin.POST("/tipo-documento", tipoDocumentoHandler.Create)
			admin.GET("/tipo-documento", tipoDocumentoHandler.List)
			admin.PUT("/tipo-documento/:id", tipoDocumentoHandler.Update)
			admin.PATCH("/tipo-documento/:id/activo", tipoDocumentoHandler.UpdateStatus)

			// =========================
			// Categoria_proveedor
			// =========================

			admin.POST("/categoria-proveedor", categoriaProveedorHandler.Create)
			admin.GET("/categoria-proveedor", categoriaProveedorHandler.List)
			admin.PUT("/categoria-proveedor/:id", categoriaProveedorHandler.Update)
			admin.PATCH("/categoria-proveedor/:id/activo", categoriaProveedorHandler.UpdateStatus)

			// =========================
			// Tipo Persona
			// =========================

			admin.POST("/tipo-persona", tipoPersonaHandler.Create)
			admin.GET("/tipo-persona", tipoPersonaHandler.List)
			admin.PUT("/tipo-persona/:id", tipoPersonaHandler.Update)
			admin.PATCH("/tipo-persona/:id/activo", tipoPersonaHandler.UpdateStatus)

			// =========================
			// Proveedor
			// =========================

			admin.POST("/proveedor", proveedorHandler.Create)
			admin.GET("/proveedor", proveedorHandler.List)
			admin.PUT("/proveedor/:id", proveedorHandler.Update)
			admin.PATCH("/proveedor/:id/activo", proveedorHandler.UpdateStatus)

			// =========================
			//	Contacto
			// =========================

			admin.POST("/contacto", contactoHandler.Create)
			admin.GET("/contacto", contactoHandler.List)
			admin.PUT("/contacto/:id", contactoHandler.Update)
			admin.PATCH("/contacto/:id/activo", contactoHandler.UpdateStatus)

			// =========================
			// Responsabilidad Fiscal
			// =========================

			admin.POST("/responsabilidad-fiscal", responsabilidadFiscalHandler.Create)
			admin.GET("/responsabilidad-fiscal", responsabilidadFiscalHandler.List)
			admin.PUT("/responsabilidad-fiscal/:id", responsabilidadFiscalHandler.Update)
			admin.PATCH("/responsabilidad-fiscal/:id/activo", responsabilidadFiscalHandler.UpdateStatus)

			// =========================
			// Receptor
			// =========================

			admin.POST("/receptor", receptorHandler.Create)
			admin.GET("/receptor", receptorHandler.List)
			admin.PUT("/receptor/:id", receptorHandler.Update)
			admin.PATCH("/receptor/:id/activo", receptorHandler.UpdateStatus)

			// =========================
			// Estado Correo
			// =========================
					
			admin.POST("/estado-correo", estadoCorreoHandler.Create)
			admin.GET("/estado-correo", estadoCorreoHandler.List)
			admin.PUT("/estado-correo/:id", estadoCorreoHandler.Update)
			admin.PATCH("/estado-correo/:id/activo", estadoCorreoHandler.UpdateStatus)
		
		}
	}

	return router
}

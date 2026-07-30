package db

// AllModels lista todos los modelos, útil para AutoMigrate:
//
//	db.AutoMigrate(models.AllModels()...)
func AllModels() []interface{} {
	return []interface{}{
		&Pais{},
		&Departamento{},
		&Municipio{},
		&Direccion{},
		&TipoDocumento{},
		&TipoPersona{},
		&CategoriaProveedor{},
		&ActividadEconomica{},
		&Proveedor{},
		&ResponsabilidadFiscal{},
		&Contacto{},
		&Rol{},
		&Usuario{},
		&Ruta{},
		&PasoRuta{},
		&Moneda{},
		&Area{},
		&Receptor{},
		&EstadoCorreo{},
		&Correo{},
		&ArchivoOrigen{},
		&Archivo{},
		&DocumentoComercial{},
		&DetalleDocumentoComercial{},
		&TipoFactura{},
		&TipoPago{},
		&MetodoPago{},
		&TipoRadicacion{},
		&CodigoQr{},
		&EstadoDocumentoRadicado{},
		&DocumentoRadicado{},
		&EstadoTarea{},
		&Tarea{},
		&Comentario{},
		&Trazabilidad{},
		&RegistroAprobacion{},
		&Notificacion{},
		&ReglaMontoRuta{},
	}
}

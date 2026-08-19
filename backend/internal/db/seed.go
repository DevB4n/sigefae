package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	// ==========================
	// Estado Correo: Recibido
	// ==========================
	estadoRecibido := EstadoCorreo{
		ID:     1,
		Nombre: "Recibido",
		Activo: true,
	}
	if err := db.Where("id = ?", estadoRecibido.ID).FirstOrCreate(&estadoRecibido).Error; err != nil {
		return err
	}

	// ==========================
	// Estados de Radicación
	// ==========================
	estadosRadicacion := []EstadoDocumentoRadicado{
		{ID: 1, Nombre: "En espera", Activo: true},
		{ID: 2, Nombre: "En proceso", Activo: true},
		{ID: 3, Nombre: "Finalizado", Activo: true},
	}
	for _, estado := range estadosRadicacion {
		if err := db.Where("id = ?", estado.ID).FirstOrCreate(&estado).Error; err != nil {
			return err
		}
	}

	// ==========================
	// Estados de Tarea
	// ==========================
	estadosTarea := []EstadoTarea{
		{ID: 1, Nombre: "Pendiente", Activo: true},
		{ID: 2, Nombre: "En Proceso", Activo: true},
		{ID: 3, Nombre: "Completada", Activo: true},
		{ID: 4, Nombre: "Devuelta", Activo: true},
	}
	for _, estado := range estadosTarea {
		if err := db.Where("id = ?", estado.ID).FirstOrCreate(&estado).Error; err != nil {
			return err
		}
	}

	// ==========================
	// Roles
	// ==========================
	superAdminRole := Rol{Nombre: "Superadministrador"}
	if err := db.Where("nombre = ?", superAdminRole.Nombre).FirstOrCreate(&superAdminRole).Error; err != nil {
		return err
	}

	aprobadorRole := Rol{Nombre: "Aprobador"}
	if err := db.Where("nombre = ?", aprobadorRole.Nombre).FirstOrCreate(&aprobadorRole).Error; err != nil {
		return err
	}

	// ==========================
	// Roles: Contabilidad y Tesorería
	// ==========================
	contabilidadRole := Rol{Nombre: "Contabilidad"}
	if err := db.Where("nombre = ?", contabilidadRole.Nombre).FirstOrCreate(&contabilidadRole).Error; err != nil {
		return err
	}

	tesoreriaRole := Rol{Nombre: "Tesorería"}
	if err := db.Where("nombre = ?", tesoreriaRole.Nombre).FirstOrCreate(&tesoreriaRole).Error; err != nil {
		return err
	}

	// ==========================
	// Usuario: Administrador
	// ==========================
	hashAdmin, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin := Usuario{
		Nombre:         "Administrador",
		Email:          "admin@sigefae.local",
		HashContrasena: string(hashAdmin),
		Cargo:          "Superadministrador",
		RolID:          superAdminRole.ID,
	}
	if err := db.Where("email = ?", admin.Email).FirstOrCreate(&admin).Error; err != nil {
		return err
	}

	// ==========================
	// Usuario: Juan Aprobador
	// ==========================
	hashAprobador, err := bcrypt.GenerateFromPassword([]byte("aprobador123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	aprobador := Usuario{
		Nombre:         "Juan Aprobador",
		Email:          "aprobador@sigefae.local",
		HashContrasena: string(hashAprobador),
		Cargo:          "Aprobador de Documentos",
		RolID:          aprobadorRole.ID,
	}
	if err := db.Where("email = ?", aprobador.Email).FirstOrCreate(&aprobador).Error; err != nil {
		return err
	}

	// ==========================
	// Usuario: Prueba Contabilidad
	// ==========================
	hashConta, _ := bcrypt.GenerateFromPassword([]byte("conta123"), bcrypt.DefaultCost)
	conta := Usuario{
		Nombre:         "Prueba Contabilidad",
		Email:          "contabilidad@sigefae.local",
		HashContrasena: string(hashConta),
		Cargo:          "Contador Principal",
		RolID:          contabilidadRole.ID,
	}
	db.Where("email = ?", conta.Email).FirstOrCreate(&conta)

	// ==========================
	// Usuario: Prueba Tesorería
	// ==========================
	hashTeso, _ := bcrypt.GenerateFromPassword([]byte("teso123"), bcrypt.DefaultCost)
	teso := Usuario{
		Nombre:         "Prueba Tesorería",
		Email:          "tesoreria@sigefae.local",
		HashContrasena: string(hashTeso),
		Cargo:          "Tesorero General",
		RolID:          tesoreriaRole.ID,
	}
	db.Where("email = ?", teso.Email).FirstOrCreate(&teso)

	// ==========================
	// Área General
	// ==========================
	areaGeneral := Area{ID: 1, Nombre: "General", Activo: true}
	if err := db.Where("id = ?", areaGeneral.ID).FirstOrCreate(&areaGeneral).Error; err != nil {
		return err
	}

	// ==========================
	// Moneda COP
	// ==========================
	monedaCOP := Moneda{ID: 1, Nombre: "Peso Colombiano", Codigo: "COP", Activo: true}
	if err := db.Where("id = ?", monedaCOP.ID).FirstOrCreate(&monedaCOP).Error; err != nil {
		return err
	}

	// ==========================
	// Origen de Archivo: Sistema  ← NUEVO
	// ==========================
	origenSistema := ArchivoOrigen{
		ID:     1,
		Nombre: "Sistema",
		Activo: true,
	}
	if err := db.Where("id = ?", origenSistema.ID).FirstOrCreate(&origenSistema).Error; err != nil {
		return err
	}
	// ==========================
	// Normas de Reparto
	// ==========================
	if err := SeedNormasReparto(db); err != nil {
		return err
	}

	return nil
}

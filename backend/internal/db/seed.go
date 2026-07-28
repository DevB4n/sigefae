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

	if err := db.
		Where("id = ?", estadoRecibido.ID).
		FirstOrCreate(&estadoRecibido).Error; err != nil {
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
		if err := db.
			Where("id = ?", estado.ID).
			FirstOrCreate(&estado).Error; err != nil {
			return err
		}
	}

	// ==========================
	// Roles
	// ==========================

	superAdminRole := Rol{
		Nombre: "Superadministrador",
	}

	if err := db.
		Where("nombre = ?", superAdminRole.Nombre).
		FirstOrCreate(&superAdminRole).Error; err != nil {
		return err
	}

	aprobadorRole := Rol{
		Nombre: "Aprobador",
	}

	if err := db.
		Where("nombre = ?", aprobadorRole.Nombre).
		FirstOrCreate(&aprobadorRole).Error; err != nil {
		return err
	}

	// ==========================
	// Usuario: Administrador
	// ==========================

	hashAdmin, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

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

	if err := db.
		Where("email = ?", admin.Email).
		FirstOrCreate(&admin).Error; err != nil {
		return err
	}

	// ==========================
	// Usuario: Aprobador
	// ==========================

	hashAprobador, err := bcrypt.GenerateFromPassword(
		[]byte("aprobador123"),
		bcrypt.DefaultCost,
	)

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

	if err := db.
		Where("email = ?", aprobador.Email).
		FirstOrCreate(&aprobador).Error; err != nil {
		return err
	}

	return nil
}
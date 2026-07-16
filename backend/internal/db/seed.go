package db

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) error {

	// ==========================
	// Rol Superadministrador
	// ==========================

	superAdminRole := Rol{
		Nombre: "Superadministrador",
	}

	if err := db.
		Where("nombre = ?", superAdminRole.Nombre).
		FirstOrCreate(&superAdminRole).Error; err != nil {
		return err
	}

	// ==========================
	// Contraseña
	// ==========================

	hash, err := bcrypt.GenerateFromPassword(
		[]byte("admin123"),
		bcrypt.DefaultCost,
	)

	if err != nil {
		return err
	}

	// ==========================
	// Usuario inicial
	// ==========================

	admin := Usuario{
		Nombre:     "Administrador",
		Email:      "admin@sigefae.local",
		HashContrasena: string(hash),
		Cargo:      "Superadministrador",
		RolID:      superAdminRole.ID,
	}

	if err := db.
		Where("email = ?", admin.Email).
		FirstOrCreate(&admin).Error; err != nil {
		return err
	}

	return nil
}
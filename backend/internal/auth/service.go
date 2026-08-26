package auth

import (
	"errors"

	"sigefae/internal/db"

	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {

	return &Service{
		db: database,
	}
}
func (s *Service) Login(email, password string) (*db.Usuario, string, error) {

	var user db.Usuario

	err := s.db.
		Preload("Rol").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("correo o contraseña incorrectos")
		}

		return nil, "", err
	}

	if !user.Activo {
		return nil, "", errors.New("usuario inactivo")
	}

	if err := Check(user.HashContrasena, password); err != nil {
		return nil, "", errors.New("correo o contraseña incorrectos")
	}

	token, err := GenerateToken(
		user.ID,
		user.RolID,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

func (s *Service) SSOLogin(email string) (*db.Usuario, string, error) {
	var user db.Usuario

	err := s.db.
		Preload("Rol").
		Where("email = ?", email).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("usuario no encontrado")
		}
		return nil, "", err
	}

	if !user.Activo {
		return nil, "", errors.New("usuario inactivo")
	}

	token, err := GenerateToken(
		user.ID,
		user.RolID,
	)

	if err != nil {
		return nil, "", err
	}

	return &user, token, nil
}

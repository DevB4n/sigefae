package user

import (
	"errors"

	"gorm.io/gorm"

	"sigefae/internal/auth"
	"sigefae/internal/db"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {

	return &Service{
		db: database,
	}
}

func (s *Service) Create(req CreateRequest) (*Response, error) {

	var existing db.Usuario

	err := s.db.
		Where("email = ?", req.Email).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el correo ya existe")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hash, err := auth.Hash(req.Contrasena)

	if err != nil {
		return nil, err
	}

	user := db.Usuario{
		Nombre:         req.Nombre,
		Email:          req.Email,
		HashContrasena: hash,
		Cargo:          req.Cargo,
		RolID:          req.RolID,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	if err := s.db.Preload("Rol").First(&user, user.ID).Error; err != nil {
		return nil, err
	}

	return &Response{
		ID:     user.ID,
		Nombre: user.Nombre,
		Email:  user.Email,
		Cargo:  user.Cargo,
		Rol:    user.Rol.Nombre,
	}, nil
}
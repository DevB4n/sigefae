package role

import (
	"errors"

	"gorm.io/gorm"

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

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Rol

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el rol ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear rol
	// ==========================

	role := db.Rol{
		Nombre: req.Nombre,
		Activo: true,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, err
	}

	response := toResponse(role)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var roles []db.Rol

	if err := s.db.
		Order("nombre ASC").
		Find(&roles).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(roles))

	for _, role := range roles {
		response = append(response, toResponse(role))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var role db.Rol

	err := s.db.
		First(&role, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("rol no encontrado")
	}

	if err != nil {
		return err
	}

	role.Activo = activo

	return s.db.Save(&role).Error
}
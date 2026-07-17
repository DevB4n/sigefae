package ruta

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

	var existing db.Ruta

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la ruta ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	var area db.Area

	err = s.db.
		First(&area, req.AreaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el área no existe")
	}

	if err != nil {
		return nil, err
	}

	ruta := db.Ruta{
		Nombre:  req.Nombre,
		Version: 1,
		AreaID:  req.AreaID,
		Activo:  true,
	}

	if err := s.db.Create(&ruta).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Area").
		First(&ruta, ruta.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(ruta)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var rutas []db.Ruta

	if err := s.db.
		Preload("Area").
		Order("nombre ASC").
		Find(&rutas).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(rutas))

	for _, ruta := range rutas {
		response = append(response, toResponse(ruta))
	}

	return response, nil
}
func (s *Service) UpdateStatus(id uint, activo bool) error {

	var ruta db.Ruta

	err := s.db.
		First(&ruta, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("ruta no encontrada")
	}

	if err != nil {
		return err
	}

	ruta.Activo = activo

	return s.db.Save(&ruta).Error
}
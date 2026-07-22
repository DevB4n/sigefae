package actividad_economica

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

func (s *Service) Create(req CreateDTO) (*Response, error) {

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.ActividadEconomica

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la actividad económica ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	err = s.db.
		Where("codigo = ?", req.Codigo).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el código ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	actividad := db.ActividadEconomica{
		Nombre: req.Nombre,
		Codigo: req.Codigo,
		Activo: true,
	}

	if err := s.db.Create(&actividad).Error; err != nil {
		return nil, err
	}

	response := toResponse(actividad)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var actividades []db.ActividadEconomica

	if err := s.db.
		Order("nombre ASC").
		Find(&actividades).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(actividades))

	for _, actividad := range actividades {
		response = append(response, toResponse(actividad))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var actividad db.ActividadEconomica

	err := s.db.First(&actividad, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("actividad económica no encontrada")
	}

	if err != nil {
		return err
	}

	actividad.Activo = activo

	return s.db.Save(&actividad).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var actividad db.ActividadEconomica

	err := s.db.First(&actividad, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("actividad económica no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.ActividadEconomica

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la actividad económica ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	err = s.db.
		Where("codigo = ? AND id <> ?", req.Codigo, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el código ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&actividad).Updates(map[string]interface{}{
		"nombre": req.Nombre,
		"codigo": req.Codigo,
	}).Error; err != nil {
		return nil, err
	}

	if err := s.db.First(&actividad, actividad.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(actividad)

	return &response, nil
}

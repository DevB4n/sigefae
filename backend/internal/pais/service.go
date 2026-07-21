package pais

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

	var existing db.Pais

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el país ya existe")
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
		return nil, errors.New("el código del país ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	pais := db.Pais{
		Nombre: req.Nombre,
		Codigo: req.Codigo,
		Activo: true,
	}

	if err := s.db.Create(&pais).Error; err != nil {
		return nil, err
	}

	response := toResponse(pais)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var paises []db.Pais

	if err := s.db.
		Order("nombre ASC").
		Find(&paises).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(paises))

	for _, pais := range paises {
		response = append(response, toResponse(pais))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var pais db.Pais

	err := s.db.
		First(&pais, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("país no encontrado")
	}

	if err != nil {
		return err
	}

	pais.Activo = activo

	return s.db.Save(&pais).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var pais db.Pais

	err := s.db.
		First(&pais, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("país no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Pais

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el país ya existe")
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
		return nil, errors.New("el código del país ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&pais).Updates(map[string]interface{}{
		"nombre": req.Nombre,
		"codigo": req.Codigo,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		First(&pais, pais.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(pais)

	return &response, nil
}
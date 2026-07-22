package departamento

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
	// Validar País
	// ==========================

	var pais db.Pais

	err := s.db.
		First(&pais, req.PaisID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el país no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Departamento

	err = s.db.
		Where("nombre = ? AND pais_id = ?", req.Nombre, req.PaisID).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el departamento ya existe para este país")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	departamento := db.Departamento{
		Nombre: req.Nombre,
		PaisID: req.PaisID,
		Activo: true,
	}

	if err := s.db.Create(&departamento).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Pais").
		First(&departamento, departamento.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(departamento)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var departamentos []db.Departamento

	if err := s.db.
		Preload("Pais").
		Order("nombre ASC").
		Find(&departamentos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(departamentos))

	for _, departamento := range departamentos {
		response = append(response, toResponse(departamento))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var departamento db.Departamento

	err := s.db.
		First(&departamento, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("departamento no encontrado")
	}

	if err != nil {
		return err
	}

	departamento.Activo = activo

	return s.db.Save(&departamento).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var departamento db.Departamento

	// IMPORTANTE:
	// Sin Preload para evitar el problema de las FK.

	err := s.db.
		First(&departamento, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("departamento no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar País
	// ==========================

	var pais db.Pais

	err = s.db.
		First(&pais, req.PaisID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el país no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Departamento

	err = s.db.
		Where(
			"nombre = ? AND pais_id = ? AND id <> ?",
			req.Nombre,
			req.PaisID,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el departamento ya existe para este país")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&departamento).Updates(map[string]interface{}{
		"nombre":  req.Nombre,
		"pais_id": req.PaisID,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("Pais").
		First(&departamento, departamento.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(departamento)

	return &response, nil
}

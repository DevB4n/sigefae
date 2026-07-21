package municipio

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
	// Validar Departamento
	// ==========================

	var departamento db.Departamento

	err := s.db.
		First(&departamento, req.DepartamentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el departamento no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Municipio

	err = s.db.
		Where("nombre = ? AND departamento_id = ?", req.Nombre, req.DepartamentoID).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el municipio ya existe para este departamento")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	municipio := db.Municipio{
		Nombre:         req.Nombre,
		DepartamentoID: req.DepartamentoID,
		Activo:         true,
	}

	if err := s.db.Create(&municipio).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Departamento").
		First(&municipio, municipio.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(municipio)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var municipios []db.Municipio

	if err := s.db.
		Preload("Departamento").
		Order("nombre ASC").
		Find(&municipios).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(municipios))

	for _, municipio := range municipios {
		response = append(response, toResponse(municipio))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var municipio db.Municipio

	err := s.db.
		First(&municipio, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("municipio no encontrado")
	}

	if err != nil {
		return err
	}

	municipio.Activo = activo

	return s.db.Save(&municipio).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var municipio db.Municipio

	// ==========================
	// Buscar municipio
	// SIN PRELOAD
	// ==========================

	err := s.db.
		First(&municipio, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("municipio no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Departamento
	// ==========================

	var departamento db.Departamento

	err = s.db.
		First(&departamento, req.DepartamentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el departamento no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar nombre repetido
	// ==========================

	var existing db.Municipio

	err = s.db.
		Where(
			"nombre = ? AND departamento_id = ? AND id <> ?",
			req.Nombre,
			req.DepartamentoID,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el municipio ya existe para este departamento")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&municipio).Updates(map[string]interface{}{
		"nombre":          req.Nombre,
		"departamento_id": req.DepartamentoID,
	}).Error; err != nil {

		return nil, err
	}

	// ==========================
	// Recargar relaciones
	// ==========================

	if err := s.db.
		Preload("Departamento").
		First(&municipio, municipio.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(municipio)

	return &response, nil
}
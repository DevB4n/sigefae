package direccion

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
	// Validar municipio
	// ==========================

	var municipio db.Municipio

	err := s.db.
		First(&municipio, req.IDMunicipio).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el municipio no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar duplicado
	// ==========================

	var existing db.Direccion

	err = s.db.
		Where(
			"nombre = ? AND linea_1 = ? AND id_municipio = ?",
			req.Nombre,
			req.Linea1,
			req.IDMunicipio,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la dirección ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	direccion := db.Direccion{
		Nombre:       req.Nombre,
		Linea1:       req.Linea1,
		Linea2:       req.Linea2,
		CodigoPostal: req.CodigoPostal,
		IDMunicipio:  req.IDMunicipio,
		Activo:       true,
	}

	if err := s.db.Create(&direccion).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Municipio").
		First(&direccion, direccion.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(direccion)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var direcciones []db.Direccion

	if err := s.db.
		Preload("Municipio").
		Order("nombre ASC").
		Find(&direcciones).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(direcciones))

	for _, direccion := range direcciones {
		response = append(response, toResponse(direccion))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var direccion db.Direccion

	err := s.db.First(&direccion, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("dirección no encontrada")
	}

	if err != nil {
		return err
	}

	direccion.Activo = activo

	return s.db.Save(&direccion).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var direccion db.Direccion

	err := s.db.
		First(&direccion, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("dirección no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar municipio
	// ==========================

	var municipio db.Municipio

	err = s.db.First(&municipio, req.IDMunicipio).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el municipio no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar duplicado
	// ==========================

	var existing db.Direccion

	err = s.db.
		Where(
			"nombre = ? AND linea_1 = ? AND id_municipio = ? AND id <> ?",
			req.Nombre,
			req.Linea1,
			req.IDMunicipio,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la dirección ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&direccion).Updates(map[string]interface{}{
		"nombre":        req.Nombre,
		"linea_1":       req.Linea1,
		"linea_2":       req.Linea2,
		"codigo_postal": req.CodigoPostal,
		"id_municipio":  req.IDMunicipio,
	}).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Municipio").
		First(&direccion, direccion.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(direccion)

	return &response, nil
}

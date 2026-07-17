package tipo_factura

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

	var existing db.TipoFactura

	err := s.db.
		Where("nombre = ? AND area_id = ?", req.Nombre, req.AreaID).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el tipo de factura ya existe para esta área")
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

	tipoFactura := db.TipoFactura{
		AreaID: req.AreaID,
		Nombre: req.Nombre,
		Activo: true,
	}

	if err := s.db.Create(&tipoFactura).Error; err != nil {
		return nil, err
	}

	response := toResponse(tipoFactura)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var tiposFactura []db.TipoFactura

	if err := s.db.
		Preload("Area").
		Order("nombre ASC").
		Find(&tiposFactura).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(tiposFactura))

	for _, tipoFactura := range tiposFactura {
		response = append(response, toResponse(tipoFactura))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var tipoFactura db.TipoFactura

	err := s.db.
		First(&tipoFactura, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("tipo de factura no encontrado")
	}

	if err != nil {
		return err
	}

	tipoFactura.Activo = activo

	return s.db.Save(&tipoFactura).Error
}
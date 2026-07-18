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

	var area db.Area

	err := s.db.First(&area, req.AreaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el área no existe")
	}

	if err != nil {
		return nil, err
	}

	var existing db.TipoFactura

	err = s.db.
		Where("nombre = ? AND area_id = ?", req.Nombre, req.AreaID).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el tipo de factura ya existe para esta área")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
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

	if err := s.db.
		Preload("Area").
		First(&tipoFactura, tipoFactura.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(tipoFactura)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var tipos []db.TipoFactura

	if err := s.db.
		Preload("Area").
		Order("nombre ASC").
		Find(&tipos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(tipos))

	for _, tipo := range tipos {
		response = append(response, toResponse(tipo))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var tipo db.TipoFactura

	err := s.db.First(&tipo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("tipo de factura no encontrado")
	}

	if err != nil {
		return err
	}

	tipo.Activo = activo

	return s.db.Save(&tipo).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var tipo db.TipoFactura

	err := s.db.First(&tipo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de factura no encontrado")
	}

	if err != nil {
		return nil, err
	}

	var area db.Area

	err = s.db.First(&area, req.AreaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el área no existe")
	}

	if err != nil {
		return nil, err
	}

	var existing db.TipoFactura

	err = s.db.
		Where(
			"nombre = ? AND area_id = ? AND id <> ?",
			req.Nombre,
			req.AreaID,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el tipo de factura ya existe para esta área")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.db.
		Model(&db.TipoFactura{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"nombre":  req.Nombre,
			"area_id": req.AreaID,
		}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("Area").
		First(&tipo, id).Error; err != nil {

		return nil, err
	}

	response := toResponse(tipo)

	return &response, nil
}
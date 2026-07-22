package detalle_documento_comercial

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

func (s *Service) Create(dto CreateDTO) (*Response, error) {

	// ==========================
	// Validar documento comercial
	// ==========================

	var documento db.DocumentoComercial

	err := s.db.First(&documento, dto.DocumentoComercialID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	detalle := db.DetalleDocumentoComercial{
		DocumentoComercialID: dto.DocumentoComercialID,
		Descripcion:          dto.Descripcion,
		ValorUnit:            dto.ValorUnit,
		IvaUnit:              dto.IvaUnit,
		Cantidad:             dto.Cantidad,
		Total:                dto.Total,
		Activo:               true,
	}

	if err := s.db.Create(&detalle).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("DocumentoComercial").
		First(&detalle, detalle.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(detalle)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var detalles []db.DetalleDocumentoComercial

	if err := s.db.
		Preload("DocumentoComercial").
		Order("id ASC").
		Find(&detalles).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(detalles))

	for _, detalle := range detalles {
		response = append(response, toResponse(detalle))
	}

	return response, nil
}

func (s *Service) Update(id uint, dto UpdateDTO) (*Response, error) {

	// ==========================
	// Validar detalle
	// ==========================

	var detalle db.DetalleDocumentoComercial

	err := s.db.First(&detalle, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el detalle del documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar documento comercial
	// ==========================

	var documento db.DocumentoComercial

	err = s.db.First(&documento, dto.DocumentoComercialID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&detalle).Updates(map[string]interface{}{
		"documento_comercial_id": dto.DocumentoComercialID,
		"descripcion":            dto.Descripcion,
		"valor_unitario":         dto.ValorUnit,
		"iva_unitario":           dto.IvaUnit,
		"cantidad":               dto.Cantidad,
		"total":                  dto.Total,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("DocumentoComercial").
		First(&detalle, id).Error; err != nil {

		return nil, err
	}

	response := toResponse(detalle)

	return &response, nil
}
func (s *Service) Delete(id uint) error {

	var detalle db.DetalleDocumentoComercial

	err := s.db.First(&detalle, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("detalle del documento comercial no encontrado")
	}

	if err != nil {
		return err
	}

	// Alternar estado
	detalle.Activo = !detalle.Activo

	return s.db.Save(&detalle).Error
}

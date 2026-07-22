package responsabilidad_fiscal

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
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err := s.db.First(&proveedor, req.IDProveedor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("proveedor no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	var existing db.ResponsabilidadFiscal

	err = s.db.
		Where(
			"id_proveedor = ? AND codigo = ?",
			req.IDProveedor,
			req.Codigo,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la responsabilidad fiscal ya existe para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	responsabilidad := db.ResponsabilidadFiscal{
		IDProveedor: req.IDProveedor,
		Codigo:      req.Codigo,
		Activo:      true,
	}

	if err := s.db.Create(&responsabilidad).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Proveedor").
		First(&responsabilidad, responsabilidad.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(responsabilidad)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var responsabilidades []db.ResponsabilidadFiscal

	if err := s.db.
		Preload("Proveedor").
		Order("codigo ASC").
		Find(&responsabilidades).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(responsabilidades))

	for _, responsabilidad := range responsabilidades {
		response = append(
			response,
			toResponse(responsabilidad),
		)
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var responsabilidad db.ResponsabilidadFiscal

	err := s.db.First(&responsabilidad, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("responsabilidad fiscal no encontrada")
	}

	if err != nil {
		return err
	}

	responsabilidad.Activo = activo

	return s.db.Save(&responsabilidad).Error
}
func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var responsabilidad db.ResponsabilidadFiscal

	err := s.db.First(&responsabilidad, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("responsabilidad fiscal no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err = s.db.First(&proveedor, req.IDProveedor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("proveedor no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar código repetido
	// ==========================

	var existing db.ResponsabilidadFiscal

	err = s.db.
		Where(
			"id_proveedor = ? AND codigo = ? AND id <> ?",
			req.IDProveedor,
			req.Codigo,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("la responsabilidad fiscal ya existe para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&responsabilidad).Updates(map[string]interface{}{
		"id_proveedor": req.IDProveedor,
		"codigo":       req.Codigo,
	}).Error; err != nil {

		return nil, err
	}

	// ==========================
	// Recargar
	// ==========================

	if err := s.db.
		Preload("Proveedor").
		First(&responsabilidad, responsabilidad.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(responsabilidad)

	return &response, nil
}

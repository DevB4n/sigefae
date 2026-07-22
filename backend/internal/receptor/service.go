package receptor

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

	var tipoDocumento db.TipoDocumento

	err := s.db.First(&tipoDocumento, req.TipoDocumentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de documento no encontrado")
	}

	if err != nil {
		return nil, err
	}

	var existing db.Receptor

	err = s.db.
		Where("numero_documento = ?", req.NumeroDocumento).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de documento ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	receptor := db.Receptor{
		Nombre:          req.Nombre,
		TipoDocumentoID: req.TipoDocumentoID,
		NumeroDocumento: req.NumeroDocumento,
		Activo:          true,
	}

	if err := s.db.Create(&receptor).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("TipoDocumento").
		First(&receptor, receptor.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(receptor)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var receptores []db.Receptor

	if err := s.db.
		Preload("TipoDocumento").
		Order("nombre ASC").
		Find(&receptores).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(receptores))

	for _, receptor := range receptores {
		response = append(response, toResponse(receptor))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var receptor db.Receptor

	err := s.db.First(&receptor, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("receptor no encontrado")
	}

	if err != nil {
		return err
	}

	receptor.Activo = activo

	return s.db.Save(&receptor).Error
}
func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var receptor db.Receptor

	err := s.db.First(&receptor, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("receptor no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Tipo Documento
	// ==========================

	var tipoDocumento db.TipoDocumento

	err = s.db.First(&tipoDocumento, req.TipoDocumentoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tipo de documento no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Documento Repetido
	// ==========================

	var existing db.Receptor

	err = s.db.
		Where("numero_documento = ? AND id <> ?", req.NumeroDocumento, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de documento ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&receptor).Updates(map[string]interface{}{
		"nombre":            req.Nombre,
		"tipo_documento_id": req.TipoDocumentoID,
		"numero_documento":  req.NumeroDocumento,
	}).Error; err != nil {

		return nil, err
	}

	// ==========================
	// Recargar con relaciones
	// ==========================

	if err := s.db.
		Preload("TipoDocumento").
		First(&receptor, receptor.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(receptor)

	return &response, nil
}

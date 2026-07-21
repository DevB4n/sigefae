package estado_documento_radicado

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

	var existing db.EstadoDocumentoRadicado

	err := s.db.
		Where("nombre = ?", req.Nombre).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el estado del documento radicado ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	estado := db.EstadoDocumentoRadicado{
		Nombre: req.Nombre,
		Activo: true,
	}

	if err := s.db.Create(&estado).Error; err != nil {
		return nil, err
	}

	response := toResponse(estado)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var estados []db.EstadoDocumentoRadicado

	if err := s.db.
		Order("nombre ASC").
		Find(&estados).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(estados))

	for _, estado := range estados {
		response = append(response, toResponse(estado))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var estado db.EstadoDocumentoRadicado

	err := s.db.First(&estado, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("estado del documento radicado no encontrado")
	}

	if err != nil {
		return err
	}

	estado.Activo = activo

	return s.db.Save(&estado).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var estado db.EstadoDocumentoRadicado

	err := s.db.First(&estado, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("estado del documento radicado no encontrado")
	}

	if err != nil {
		return nil, err
	}

	var existing db.EstadoDocumentoRadicado

	err = s.db.
		Where("nombre = ? AND id <> ?", req.Nombre, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el estado del documento radicado ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.db.Model(&estado).Updates(map[string]interface{}{
		"nombre": req.Nombre,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&estado, estado.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(estado)

	return &response, nil
}
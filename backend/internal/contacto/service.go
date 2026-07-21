package contacto

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
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err := s.db.First(&proveedor, req.ProveedorID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el proveedor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar correo repetido
	// ==========================

	var existing db.Contacto

	err = s.db.
		Where(
			"proveedor_id = ? AND correo = ?",
			req.ProveedorID,
			req.Correo,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un contacto con ese correo para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	contacto := db.Contacto{
		ProveedorID: req.ProveedorID,
		Nombre:      req.Nombre,
		Cargo:       req.Cargo,
		Telefono:    req.Telefono,
		Correo:      req.Correo,
		Activo:      true,
	}

	if err := s.db.Create(&contacto).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Proveedor").
		First(&contacto, contacto.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(contacto)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var contactos []db.Contacto

	if err := s.db.
		Preload("Proveedor").
		Order("nombre ASC").
		Find(&contactos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(contactos))

	for _, contacto := range contactos {
		response = append(response, toResponse(contacto))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, activo bool) error {

	var contacto db.Contacto

	err := s.db.First(&contacto, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("contacto no encontrado")
	}

	if err != nil {
		return err
	}

	contacto.Activo = activo

	return s.db.Save(&contacto).Error
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var contacto db.Contacto

	err := s.db.First(&contacto, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("contacto no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err = s.db.First(&proveedor, req.ProveedorID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el proveedor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar correo repetido
	// ==========================

	var existing db.Contacto

	err = s.db.
		Where(
			"proveedor_id = ? AND correo = ? AND id <> ?",
			req.ProveedorID,
			req.Correo,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un contacto con ese correo para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&contacto).Updates(map[string]interface{}{
		"proveedor_id": req.ProveedorID,
		"nombre":       req.Nombre,
		"cargo":        req.Cargo,
		"telefono":     req.Telefono,
		"correo":       req.Correo,
	}).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Proveedor").
		First(&contacto, contacto.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(contacto)

	return &response, nil
}
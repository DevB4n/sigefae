package archivo

import (
	"errors"
	"os"

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
	// Validar Documento Radicado
	// ==========================

	var documento db.DocumentoRadicado

	err := s.db.First(
		&documento,
		req.DocumentoRadicadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("documento radicado no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Origen
	// ==========================

	var origen db.ArchivoOrigen

	err = s.db.First(
		&origen,
		req.OrigenID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("origen no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta repetida
	// ==========================

	var existing db.Archivo

	err = s.db.
		Where("ruta = ?", req.Ruta).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un archivo con esa ruta")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// ==========================
	// Crear
	// ==========================

	archivo := db.Archivo{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		Nombre:              req.Nombre,
		Extension:           req.Extension,
		Peso:                req.Peso,
		Ruta:                req.Ruta,
		OrigenID:            req.OrigenID,
	}

	if err := s.db.Create(&archivo).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var archivos []db.Archivo

	if err := s.db.
		Order("created_at DESC").
		Find(&archivos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(archivos))

	for _, archivo := range archivos {
		response = append(response, toResponse(archivo))
	}

	return response, nil
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var archivo db.Archivo

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("archivo no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Origen
	// ==========================

	var origen db.ArchivoOrigen

	err = s.db.First(
		&origen,
		req.OrigenID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("origen no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta repetida
	// ==========================

	var existing db.Archivo

	err = s.db.
		Where("ruta = ? AND id <> ?", req.Ruta, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("ya existe un archivo con esa ruta")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&archivo).Updates(map[string]interface{}{
		"nombre":    req.Nombre,
		"extension": req.Extension,
		"peso":      req.Peso,
		"ruta":      req.Ruta,
		"origen_id": req.OrigenID,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&archivo, archivo.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(archivo)

	return &response, nil
}

func (s *Service) Delete(id uint) error {

	var archivo db.Archivo

	err := s.db.First(&archivo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("archivo no encontrado")
	}

	if err != nil {
		return err
	}

	// ==========================
	// Eliminar archivo físico
	// ==========================

	if archivo.Ruta != "" {

		if err := os.Remove(archivo.Ruta); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	// ==========================
	// Eliminar registro
	// ==========================

	return s.db.Delete(&archivo).Error
}

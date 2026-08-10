package tarea

import (
	"errors"
	"time"

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
	// Validar Usuario
	// ==========================

	var usuario db.Usuario

	err = s.db.First(
		&usuario,
		req.UsuarioAsignadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("usuario no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Estado
	// ==========================

	var estado db.EstadoTarea

	err = s.db.First(
		&estado,
		req.EstadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("estado de tarea no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	tarea := db.Tarea{
		DocumentoRadicadoID: req.DocumentoRadicadoID,
		UsuarioAsignadoID:   req.UsuarioAsignadoID,
		EstadoID:            req.EstadoID,
		Descripcion:         req.Descripcion,
		FechaAsignacion:     time.Now(),
		FechaLimite:         req.FechaLimite,
	}

	if err := s.db.Create(&tarea).Error; err != nil {
		return nil, err
	}

	response := toResponse(tarea)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var tareas []db.Tarea

	if err := s.db.
		Order("created_at DESC").
		Find(&tareas).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(tareas))

	for _, tarea := range tareas {
		response = append(response, toResponse(tarea))
	}

	return response, nil
}

func (s *Service) Update(id uint, req UpdateDTO) (*Response, error) {

	var tarea db.Tarea

	err := s.db.First(&tarea, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("tarea no encontrada")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Usuario
	// ==========================

	var usuario db.Usuario

	err = s.db.First(
		&usuario,
		req.UsuarioAsignadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("usuario no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar Estado
	// ==========================

	var estado db.EstadoTarea

	err = s.db.First(
		&estado,
		req.EstadoID,
	).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("estado de tarea no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&tarea).Updates(map[string]any{
		"usuario_asignado_id": req.UsuarioAsignadoID,
		"estado_id":           req.EstadoID,
		"descripcion":         req.Descripcion,
		"fecha_inicio":        req.FechaInicio,
		"fecha_limite":        req.FechaLimite,
		"fecha_finalizacion":  req.FechaFinalizacion,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.First(&tarea, tarea.ID).Error; err != nil {
		return nil, err
	}

	response := toResponse(tarea)

	return &response, nil
}

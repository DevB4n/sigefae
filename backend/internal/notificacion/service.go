package notificacion

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

func (s *Service) Create(dto CreateDTO) (*Response, error) {

	// ==========================
	// Validar usuario
	// ==========================

	var usuario db.Usuario

	err := s.db.First(&usuario, dto.UsuarioID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el usuario no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar documento radicado (opcional)
	// ==========================

	if dto.DocumentoRadicadoID != nil {

		var documento db.DocumentoRadicado

		err = s.db.First(&documento, *dto.DocumentoRadicadoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el documento radicado no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar estado
	// ==========================

	switch dto.Estado {
	case "Pendiente", "Enviada", "Leida", "Expirada":
	default:
		return nil, errors.New("estado de notificación inválido")
	}

	// ==========================
	// Validar tipo
	// ==========================

	switch dto.Tipo {
	case "Recordatorio", "Asignacion", "Sistema", "Rechazo", "Finalizado", "Devolucion":
	default:
		return nil, errors.New("tipo de notificación inválido")
	}

	// ==========================
	// Crear
	// ==========================

	notificacion := db.Notificacion{
		UsuarioID:           dto.UsuarioID,
		DocumentoRadicadoID: dto.DocumentoRadicadoID,
		Mensaje:             dto.Mensaje,
		Estado:              dto.Estado,
		Tipo:                dto.Tipo,
		FechaCreacion:       dto.FechaCreacion,
		FechaEnvio:          dto.FechaEnvio,
		FechaLectura:        dto.FechaLectura,
	}

	if err := s.db.Create(&notificacion).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Usuario").
		Preload("DocumentoRadicado").
		First(&notificacion, notificacion.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(notificacion)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var notificaciones []db.Notificacion

	if err := s.db.
		Preload("Usuario").
		Preload("DocumentoRadicado").
		Order("fecha_creacion DESC").
		Find(&notificaciones).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(notificaciones))

	for _, notificacion := range notificaciones {
		response = append(response, toResponse(notificacion))
	}

	return response, nil
}
func (s *Service) Update(id uint, dto UpdateDTO) (*Response, error) {

	// ==========================
	// Validar notificación
	// ==========================

	var notificacion db.Notificacion

	err := s.db.First(&notificacion, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("la notificación no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar usuario
	// ==========================

	var usuario db.Usuario

	err = s.db.First(&usuario, dto.UsuarioID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el usuario no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar documento radicado (opcional)
	// ==========================

	if dto.DocumentoRadicadoID != nil {

		var documento db.DocumentoRadicado

		err = s.db.First(&documento, *dto.DocumentoRadicadoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el documento radicado no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar estado
	// ==========================

	switch dto.Estado {
	case "Pendiente", "Enviada", "Leida", "Expirada":
	default:
		return nil, errors.New("estado de notificación inválido")
	}

	// ==========================
	// Validar tipo
	// ==========================

	switch dto.Tipo {
	case "Recordatorio", "Asignacion", "Sistema", "Rechazo", "Finalizado", "Devolucion":
	default:
		return nil, errors.New("tipo de notificación inválido")
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&notificacion).Updates(map[string]any{
		"usuario_id":            dto.UsuarioID,
		"documento_radicado_id": dto.DocumentoRadicadoID,
		"mensaje":               dto.Mensaje,
		"estado":                dto.Estado,
		"tipo":                  dto.Tipo,
		"fecha_creacion":        dto.FechaCreacion,
		"fecha_envio":           dto.FechaEnvio,
		"fecha_lectura":         dto.FechaLectura,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("Usuario").
		Preload("DocumentoRadicado").
		First(&notificacion, id).Error; err != nil {

		return nil, err
	}

	response := toResponse(notificacion)

	return &response, nil
}
func (s *Service) ListByUsuario(usuarioID uint) ([]Response, error) {
	var notificaciones []db.Notificacion
	if err := s.db.
		Preload("Usuario").
		Preload("DocumentoRadicado").
		Where("usuario_id = ?", usuarioID).
		Order("fecha_creacion DESC").
		Find(&notificaciones).Error; err != nil {
		return nil, err
	}
	response := make([]Response, 0, len(notificaciones))
	for _, n := range notificaciones {
		response = append(response, toResponse(n))
	}
	return response, nil
}

func (s *Service) MarkAsRead(id uint) (*Response, error) {
	var notificacion db.Notificacion
	if err := s.db.First(&notificacion, id).Error; err != nil {
		return nil, err
	}
	now := time.Now()
	notificacion.Estado = "Leida"
	notificacion.FechaLectura = &now
	if err := s.db.Save(&notificacion).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Usuario").Preload("DocumentoRadicado").First(&notificacion, id)
	resp := toResponse(notificacion)
	return &resp, nil
}
func (s *Service) CreateFromEvent(dto CreateDTO) (*Response, error) {
	// Normalizar estado
	switch dto.Estado {
	case "Pendiente", "Enviada", "Leida", "Expirada":
	default:
		dto.Estado = "Pendiente"
	}
	// Normalizar tipo
	switch dto.Tipo {
	case "Recordatorio", "Asignacion", "Sistema", "Rechazo", "Finalizado", "Devolucion":
	default:
		dto.Tipo = "Sistema"
	}

	n := db.Notificacion{
		UsuarioID:           dto.UsuarioID,
		DocumentoRadicadoID: dto.DocumentoRadicadoID,
		Mensaje:             dto.Mensaje,
		Estado:              dto.Estado,
		Tipo:                dto.Tipo,
		FechaCreacion:       dto.FechaCreacion,
	}
	if err := s.db.Create(&n).Error; err != nil {
		return nil, err
	}
	s.db.Preload("Usuario").Preload("DocumentoRadicado").First(&n, n.ID)
	resp := toResponse(n)
	return &resp, nil
}

func (s *Service) Delete(id uint) error {
	var notificacion db.Notificacion
	if err := s.db.First(&notificacion, id).Error; err != nil {
		return err
	}
	return s.db.Delete(&notificacion).Error
}

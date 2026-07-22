package correo

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
	// Validar estado
	// ==========================

	var estado db.EstadoCorreo

	err := s.db.First(&estado, req.IDEstado).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el estado del correo no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ID mensaje
	// ==========================

	var existing db.Correo

	err = s.db.
		Where("id_mensaje = ?", req.IDMensaje).
		First(&existing).Error

	if err == nil {

		if existing.Activo {
			return nil, errors.New("el correo ya existe")
		}

		// ==========================
		// Reactivar
		// ==========================

		existing.Asunto = req.Asunto
		existing.De = req.De
		existing.Para = req.Para
		existing.FechaRecepcion = req.FechaRecepcion
		existing.Cuerpo = req.Cuerpo
		existing.Cc = req.Cc
		existing.Bcc = req.Bcc
		existing.ReplyTo = req.ReplyTo
		existing.IDEstado = req.IDEstado
		existing.Activo = true

		if err := s.db.Save(&existing).Error; err != nil {
			return nil, err
		}

		if err := s.db.
			Preload("EstadoCorreo").
			First(&existing, existing.ID).Error; err != nil {

			return nil, err
		}

		response := toResponse(existing)

		return &response, nil
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear
	// ==========================

	correo := db.Correo{
		Asunto:         req.Asunto,
		De:             req.De,
		Para:           req.Para,
		FechaRecepcion: req.FechaRecepcion,
		IDMensaje:      req.IDMensaje,
		Cuerpo:         req.Cuerpo,
		Cc:             req.Cc,
		Bcc:            req.Bcc,
		ReplyTo:        req.ReplyTo,
		IDEstado:       req.IDEstado,
		Activo:         true,
	}

	if err := s.db.Create(&correo).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("EstadoCorreo").
		First(&correo, correo.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(correo)

	return &response, nil
}

func (s *Service) List() ([]Response, error) {

	var correos []db.Correo

	if err := s.db.
		Preload("EstadoCorreo").
		Where("activo = ?", true).
		Order("fecha_recepcion DESC").
		Find(&correos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(correos))

	for _, correo := range correos {
		response = append(response, toResponse(correo))
	}

	return response, nil
}

func (s *Service) UpdateStatus(id uint, idEstado uint) error {

	var correo db.Correo

	err := s.db.First(&correo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("correo no encontrado")
	}

	if err != nil {
		return err
	}

	var estado db.EstadoCorreo

	err = s.db.First(&estado, idEstado).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("el estado del correo no existe")
	}

	if err != nil {
		return err
	}

	correo.IDEstado = idEstado

	return s.db.Save(&correo).Error
}

func (s *Service) Update(id uint, req UpdateRequest) (*Response, error) {

	var correo db.Correo

	err := s.db.First(&correo, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("correo no encontrado")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar estado
	// ==========================

	var estado db.EstadoCorreo

	err = s.db.First(&estado, req.IDEstado).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el estado del correo no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ID mensaje
	// ==========================

	var existing db.Correo

	err = s.db.
		Where("id_mensaje = ? AND id <> ?", req.IDMensaje, id).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el id del mensaje ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&correo).Updates(map[string]interface{}{
		"asunto":          req.Asunto,
		"de":              req.De,
		"para":            req.Para,
		"fecha_recepcion": req.FechaRecepcion,
		"id_mensaje":      req.IDMensaje,
		"cuerpo":          req.Cuerpo,
		"cc":              req.Cc,
		"bcc":             req.Bcc,
		"reply_to":        req.ReplyTo,
		"id_estado":       req.IDEstado,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("EstadoCorreo").
		First(&correo, correo.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(correo)

	return &response, nil
}

func (s *Service) Delete(id uint) error {

	result := s.db.
		Model(&db.Correo{}).
		Where("id = ? AND activo = ?", id, true).
		Update("activo", false)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("correo no encontrado o ya eliminado")
	}

	return nil
}
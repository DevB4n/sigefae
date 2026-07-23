package documento_radicado

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

	var documentoComercial db.DocumentoComercial

	err := s.db.First(&documentoComercial, dto.DocumentoComercialID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar tipo de radicación
	// ==========================

	var tipoRadicacion db.TipoRadicacion

	err = s.db.First(&tipoRadicacion, dto.TipoRadicacionID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el tipo de radicación no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta
	// ==========================

	var ruta db.Ruta

	err = s.db.First(&ruta, dto.RutaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("la ruta no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar usuario actual
	// ==========================

	var usuario db.Usuario

	err = s.db.First(&usuario, dto.UsuarioActualID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el usuario no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar paso actual (opcional)
	// ==========================

	if dto.PasoActualID != nil {

		var paso db.PasoRuta

		err = s.db.First(&paso, *dto.PasoActualID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el paso actual no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar paso pendiente retorno (opcional)
	// ==========================

	if dto.PasoPendienteRetornoID != nil {

		var paso db.PasoRuta

		err = s.db.First(&paso, *dto.PasoPendienteRetornoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el paso pendiente de retorno no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar estado
	// ==========================

	var estado db.EstadoDocumentoRadicado

	err = s.db.First(&estado, dto.EstadoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el estado del documento radicado no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar código QR
	// ==========================

	var qr db.CodigoQr

	err = s.db.First(&qr, dto.QrID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el código QR no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar método de pago
	// ==========================

	var metodoPago db.MetodoPago

	err = s.db.First(&metodoPago, dto.MetodoPagoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el método de pago no existe")
	}

	if err != nil {
		return nil, err
	}
	// ==========================
	// Validar documento único
	// ==========================

	var existing db.DocumentoRadicado

	err = s.db.
		Where("documento_comercial_id = ?", dto.DocumentoComercialID).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el documento comercial ya fue radicado")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar número de radicado
	// ==========================

	err = s.db.
		Where("numero_radicado = ?", dto.NumeroRadicado).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de radicado ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Crear documento radicado
	// ==========================

	documento := db.DocumentoRadicado{
		DocumentoComercialID:   dto.DocumentoComercialID,
		TipoRadicacionID:       dto.TipoRadicacionID,
		RutaID:                 dto.RutaID,
		NumeroRadicado:         dto.NumeroRadicado,
		FechaRadicacion:        dto.FechaRadicacion,
		UsuarioActualID:        dto.UsuarioActualID,
		EstadoPosesion:         dto.EstadoPosesion,
		PasoActualID:           dto.PasoActualID,
		PasoPendienteRetornoID: dto.PasoPendienteRetornoID,
		EstadoID:               dto.EstadoID,
		UltimaActividad:        dto.UltimaActividad,
		QrID:                   dto.QrID,
		MetodoPagoID:           dto.MetodoPagoID,
	}

	if err := s.db.Create(&documento).Error; err != nil {
		return nil, err
	}

	// ==========================
	// Cargar relaciones
	// ==========================

	if err := s.db.
		Preload("DocumentoComercial").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("PasoPendienteRetorno").
		Preload("Estado").
		Preload("Qr").
		Preload("MetodoPago").
		First(&documento, documento.ID).Error; err != nil {

		return nil, err
	}
	response := toResponse(documento)
	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var documentos []db.DocumentoRadicado

	if err := s.db.
		Preload("DocumentoComercial").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("PasoPendienteRetorno").
		Preload("Estado").
		Preload("Qr").
		Preload("MetodoPago").
		Order("fecha_radicacion DESC").
		Find(&documentos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(documentos))

	for _, documento := range documentos {
		response = append(response, toResponse(documento))
	}

	return response, nil
}
func (s *Service) Update(id uint, dto UpdateDTO) (*Response, error) {

	// ==========================
	// Validar documento radicado
	// ==========================

	var documento db.DocumentoRadicado

	err := s.db.First(&documento, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento radicado no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar documento comercial
	// ==========================

	var documentoComercial db.DocumentoComercial

	err = s.db.First(&documentoComercial, dto.DocumentoComercialID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar tipo de radicación
	// ==========================

	var tipoRadicacion db.TipoRadicacion

	err = s.db.First(&tipoRadicacion, dto.TipoRadicacionID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el tipo de radicación no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar ruta
	// ==========================

	var ruta db.Ruta

	err = s.db.First(&ruta, dto.RutaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("la ruta no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar usuario actual
	// ==========================

	var usuario db.Usuario

	err = s.db.First(&usuario, dto.UsuarioActualID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el usuario no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar paso actual (opcional)
	// ==========================

	if dto.PasoActualID != nil {

		var paso db.PasoRuta

		err = s.db.First(&paso, *dto.PasoActualID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el paso actual no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar paso pendiente retorno (opcional)
	// ==========================

	if dto.PasoPendienteRetornoID != nil {

		var paso db.PasoRuta

		err = s.db.First(&paso, *dto.PasoPendienteRetornoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el paso pendiente de retorno no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar estado
	// ==========================

	var estado db.EstadoDocumentoRadicado

	err = s.db.First(&estado, dto.EstadoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el estado del documento radicado no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar código QR
	// ==========================

	var qr db.CodigoQr

	err = s.db.First(&qr, dto.QrID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el código QR no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar método de pago
	// ==========================

	var metodoPago db.MetodoPago

	err = s.db.First(&metodoPago, dto.MetodoPagoID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el método de pago no existe")
	}

	if err != nil {
		return nil, err
	}
	// ==========================
	// Validar documento comercial único
	// ==========================

	var existing db.DocumentoRadicado

	err = s.db.
		Where(
			"documento_comercial_id = ? AND id <> ?",
			dto.DocumentoComercialID,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el documento comercial ya fue radicado")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar número de radicado único
	// ==========================

	err = s.db.
		Where(
			"numero_radicado = ? AND id <> ?",
			dto.NumeroRadicado,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el número de radicado ya existe")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&documento).Updates(map[string]interface{}{
		"documento_comercial_id":    dto.DocumentoComercialID,
		"tipo_radicacion_id":        dto.TipoRadicacionID,
		"ruta_id":                   dto.RutaID,
		"numero_radicado":           dto.NumeroRadicado,
		"fecha_radicacion":          dto.FechaRadicacion,
		"usuario_actual_id":         dto.UsuarioActualID,
		"estado_posesion":           dto.EstadoPosesion,
		"paso_actual_id":            dto.PasoActualID,
		"paso_pendiente_retorno_id": dto.PasoPendienteRetornoID,
		"estado_id":                 dto.EstadoID,
		"ultima_actividad":          dto.UltimaActividad,
		"qr_id":                     dto.QrID,
		"metodo_pago_id":            dto.MetodoPagoID,
	}).Error; err != nil {

		return nil, err
	}

	// ==========================
	// Cargar relaciones
	// ==========================

	if err := s.db.
		Preload("DocumentoComercial").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("PasoPendienteRetorno").
		Preload("Estado").
		Preload("Qr").
		Preload("MetodoPago").
		First(&documento, documento.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(documento)

	return &response, nil
}

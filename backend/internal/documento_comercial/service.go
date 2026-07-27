package documento_comercial

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
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err := s.db.First(&proveedor, dto.IDProveedor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el proveedor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar receptor
	// ==========================

	var receptor db.Receptor

	err = s.db.First(&receptor, dto.IDReceptor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el receptor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar área
	// ==========================

	var area db.Area

	err = s.db.First(&area, dto.IDArea).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el área no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar moneda
	// ==========================

	var moneda db.Moneda

	err = s.db.First(&moneda, dto.MonedaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("la moneda no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar tipo factura (opcional)
	// ==========================

	if dto.TipoFacturaID != nil {

		var tipoFactura db.TipoFactura

		err = s.db.First(&tipoFactura, *dto.TipoFacturaID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el tipo de factura no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar correo (opcional)
	// ==========================

	if dto.CorreoID != nil {

		var correo db.Correo

		err = s.db.First(&correo, *dto.CorreoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el correo no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar proveedor + documento
	// ==========================

	var existing db.DocumentoComercial

	err = s.db.
		Where(
			"id_proveedor = ? AND numero_documento = ?",
			dto.IDProveedor,
			dto.NumeroDocumento,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el documento comercial ya existe para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar CUFE
	// ==========================

	if dto.Cufe != nil {

		var existingCUFE db.DocumentoComercial

		err = s.db.
			Where("cufe = ?", *dto.Cufe).
			First(&existingCUFE).Error

		if err == nil {
			return nil, errors.New("el CUFE ya existe")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// ==========================
	// Validar correo único
	// ==========================

	if dto.CorreoID != nil {

		var existingCorreo db.DocumentoComercial

		err = s.db.
			Where("correo_id = ?", *dto.CorreoID).
			First(&existingCorreo).Error

		if err == nil {
			return nil, errors.New("el correo ya está asociado a otro documento")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// ==========================
	// Crear documento
	// ==========================

	documento := db.DocumentoComercial{
		Tipo:                     dto.Tipo,
		NumeroDocumento:          dto.NumeroDocumento,
		OrdenCompra:              dto.OrdenCompra,
		IDProveedor:              dto.IDProveedor,
		IDReceptor:               dto.IDReceptor,
		IDArea:                   dto.IDArea,
		TipoFacturaID:            dto.TipoFacturaID,
		Asunto:                   dto.Asunto,
		FechaDocumento:           dto.FechaDocumento,
		FechaVencimiento:         dto.FechaVencimiento,
		MonedaID:                 dto.MonedaID,
		Subtotal:                 dto.Subtotal,
		Iva:                      dto.Iva,
		Total:                    dto.Total,
		NumeroFolios:             dto.NumeroFolios,
		OrientacionSelloRecibido: dto.OrientacionSelloRecibido,
		Cufe:                     dto.Cufe,
		CorreoID:                 dto.CorreoID,
		Activo:                   dto.Activo,
	}

	if err := s.db.Create(&documento).Error; err != nil {
		return nil, err
	}

	if err := s.db.
		Preload("Proveedor").
		Preload("Receptor").
		Preload("Area").
		Preload("TipoFactura").
		Preload("Moneda").
		Preload("Correo").
		First(&documento, documento.ID).Error; err != nil {

		return nil, err
	}

	response := toResponse(documento)

	return &response, nil
}
func (s *Service) List() ([]Response, error) {

	var documentos []db.DocumentoComercial

	if err := s.db.
		Preload("Proveedor").
		Preload("Receptor").
		Preload("Area").
		Preload("TipoFactura").
		Preload("Moneda").
		Preload("Correo").
		Order("created_at DESC").
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
	// Validar documento
	// ==========================

	var documento db.DocumentoComercial

	err := s.db.First(&documento, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el documento comercial no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar proveedor
	// ==========================

	var proveedor db.Proveedor

	err = s.db.First(&proveedor, dto.IDProveedor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el proveedor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar receptor
	// ==========================

	var receptor db.Receptor

	err = s.db.First(&receptor, dto.IDReceptor).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el receptor no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar área
	// ==========================

	var area db.Area

	err = s.db.First(&area, dto.IDArea).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("el área no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar moneda
	// ==========================

	var moneda db.Moneda

	err = s.db.First(&moneda, dto.MonedaID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("la moneda no existe")
	}

	if err != nil {
		return nil, err
	}

	// ==========================
	// Validar tipo factura
	// ==========================

	if dto.TipoFacturaID != nil {

		var tipoFactura db.TipoFactura

		err = s.db.First(&tipoFactura, *dto.TipoFacturaID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el tipo de factura no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar correo
	// ==========================

	if dto.CorreoID != nil {

		var correo db.Correo

		err = s.db.First(&correo, *dto.CorreoID).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el correo no existe")
		}

		if err != nil {
			return nil, err
		}
	}

	// ==========================
	// Validar proveedor + documento
	// ==========================

	var existing db.DocumentoComercial

	err = s.db.
		Where(
			"id_proveedor = ? AND numero_documento = ? AND id <> ?",
			dto.IDProveedor,
			dto.NumeroDocumento,
			id,
		).
		First(&existing).Error

	if err == nil {
		return nil, errors.New("el documento comercial ya existe para este proveedor")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// ==========================
	// Validar CUFE
	// ==========================

	if dto.Cufe != nil {

		var existingCUFE db.DocumentoComercial

		err = s.db.
			Where("cufe = ? AND id <> ?", *dto.Cufe, id).
			First(&existingCUFE).Error

		if err == nil {
			return nil, errors.New("el CUFE ya existe")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// ==========================
	// Validar correo único
	// ==========================

	if dto.CorreoID != nil {

		var existingCorreo db.DocumentoComercial

		err = s.db.
			Where("correo_id = ? AND id <> ?", *dto.CorreoID, id).
			First(&existingCorreo).Error

		if err == nil {
			return nil, errors.New("el correo ya está asociado a otro documento")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	// ==========================
	// Actualizar
	// ==========================

	if err := s.db.Model(&documento).Updates(map[string]interface{}{
		"tipo":                       dto.Tipo,
		"numero_documento":           dto.NumeroDocumento,
		"orden_compra":               dto.OrdenCompra,
		"id_proveedor":               dto.IDProveedor,
		"id_receptor":                dto.IDReceptor,
		"id_area":                    dto.IDArea,
		"tipo_factura_id":            dto.TipoFacturaID,
		"asunto":                     dto.Asunto,
		"fecha_documento":            dto.FechaDocumento,
		"fecha_vencimiento":          dto.FechaVencimiento,
		"moneda_id":                  dto.MonedaID,
		"subtotal":                   dto.Subtotal,
		"iva":                        dto.Iva,
		"total":                      dto.Total,
		"numero_folios":              dto.NumeroFolios,
		"orientacion_sello_recibido": dto.OrientacionSelloRecibido,
		"cufe":                       dto.Cufe,
		"correo_id":                  dto.CorreoID,
		"activo":                     dto.Activo,
	}).Error; err != nil {

		return nil, err
	}

	if err := s.db.
		Preload("Proveedor").
		Preload("Receptor").
		Preload("Area").
		Preload("TipoFactura").
		Preload("Moneda").
		Preload("Correo").
		First(&documento, id).Error; err != nil {

		return nil, err
	}

	response := toResponse(documento)

	return &response, nil
}

func (s *Service) Delete(id uint) error {

	var documento db.DocumentoComercial

	if err := s.db.First(&documento, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("documento comercial no encontrado")
		}

		return err
	}

	documento.Activo = !documento.Activo

	return s.db.Save(&documento).Error
}

func (s *Service) GetByID(id uint) (*Response, error) {

	var documento db.DocumentoComercial

	if err := s.db.
		Preload("Proveedor").
		Preload("Proveedor.TipoDocumento").
		Preload("Proveedor.Direccion").
		Preload("Proveedor.Direccion.Municipio").
		Preload("Proveedor.Direccion.Municipio.Departamento").
		Preload("Receptor").
		Preload("Area").
		Preload("TipoFactura").
		Preload("Moneda").
		Preload("Correo").
		Preload("Correo.EstadoCorreo").
		Preload("Detalles").
		First(&documento, id).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el documento comercial no existe")
		}

		return nil, err
	}

	response := toResponse(documento)

	return &response, nil
}

// ListPendientes devuelve documentos comerciales que están "Procesados" (extraídos del XML)
// pero que aún NO han sido radicados. Estos son los que la persona debe revisar.
func (s *Service) ListPendientes() ([]Response, error) {

	var documentos []db.DocumentoComercial

	// Documentos que:
	// 1. Tienen correo asociado (origen digital)
	// 2. NO tienen un documento_radicado asociado
	subQuery := s.db.Table("documento_radicado").Select("documento_comercial_id")

	if err := s.db.
		Preload("Proveedor").
		Preload("Receptor").
		Preload("Area").
		Preload("TipoFactura").
		Preload("Moneda").
		Preload("Correo").
		Preload("Correo.EstadoCorreo").
		Preload("Detalles").
		Where("id NOT IN (?)", subQuery).
		Where("activo = ?", true).
		Order("created_at DESC").
		Find(&documentos).Error; err != nil {

		return nil, err
	}

	response := make([]Response, 0, len(documentos))

	for _, documento := range documentos {
		response = append(response, toResponse(documento))
	}

	return response, nil
}

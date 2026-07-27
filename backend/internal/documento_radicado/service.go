package documento_radicado

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sigefae/internal/db"
)

type Service struct {
	db *gorm.DB
}

func New(database *gorm.DB) *Service {
	return &Service{db: database}
}

func (s *Service) Create(dto CreateDTO, usuarioID uint) (*db.DocumentoRadicado, error) {
	// ── 1. Validar que el documento comercial existe ──
	var docCom db.DocumentoComercial
	if err := s.db.First(&docCom, dto.DocumentoComercialID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el documento comercial no existe")
		}
		return nil, err
	}

	// ── 2. Validar que NO esté ya radicado ──
	var existing db.DocumentoRadicado
	if err := s.db.Where("documento_comercial_id = ?", dto.DocumentoComercialID).First(&existing).Error; err == nil {
		return nil, errors.New("el documento ya fue radicado anteriormente")
	}

	// ── 3. Validar catálogos ──
	var tipoRad db.TipoRadicacion
	if err := s.db.First(&tipoRad, dto.TipoRadicacionID).Error; err != nil {
		return nil, errors.New("el tipo de radicación no existe")
	}

	var ruta db.Ruta
	if err := s.db.First(&ruta, dto.RutaID).Error; err != nil {
		return nil, errors.New("la ruta no existe")
	}

	var metodoPago db.MetodoPago
	if err := s.db.First(&metodoPago, dto.MetodoPagoID).Error; err != nil {
		return nil, errors.New("el método de pago no existe")
	}

	// ── 4. Obtener primer paso de la ruta ──
	var primerPaso db.PasoRuta
	var primerPasoID *uint
	err := s.db.Where("ruta_id = ?", dto.RutaID).Order("id ASC").First(&primerPaso).Error
	if err == nil {
		id := primerPaso.ID
		primerPasoID = &id
	}

	// ── 5. Obtener estado inicial (el primero por ID) ──
	var estadoInicial db.EstadoDocumentoRadicado
	if err := s.db.Order("id ASC").First(&estadoInicial).Error; err != nil {
		return nil, errors.New("no hay estados de radicación configurados en el sistema")
	}

	// ── 6. Generar número de radicado ──
	numeroRadicado := dto.NumeroRadicado
	if numeroRadicado == "" {
		year := time.Now().Year()
		var count int64
		s.db.Model(&db.DocumentoRadicado{}).
			Where("YEAR(fecha_radicacion) = ?", year).
			Count(&count)
		numeroRadicado = fmt.Sprintf("RAD-%d-%05d", year, count+1)
	}

	// ── 7. Validar que el número no exista ──
	var existingNum db.DocumentoRadicado
	if err := s.db.Where("numero_radicado = ?", numeroRadicado).First(&existingNum).Error; err == nil {
		return nil, errors.New("el número de radicado ya existe")
	}

	// ── 8. Crear QR ──
	qr := db.CodigoQr{
		Url:    fmt.Sprintf("https://sigefae.com/radicado/%s", numeroRadicado),
		Activo: true,
	}
	if err := s.db.Create(&qr).Error; err != nil {
		return nil, err
	}

	// ── 9. Crear radicado ──
	radicado := db.DocumentoRadicado{
		DocumentoComercialID: dto.DocumentoComercialID,
		TipoRadicacionID:     dto.TipoRadicacionID,
		RutaID:               dto.RutaID,
		NumeroRadicado:       numeroRadicado,
		FechaRadicacion:      time.Now(),
		UsuarioActualID:      usuarioID,
		EstadoPosesion:       "Libre",
		PasoActualID:         primerPasoID,
		EstadoID:             estadoInicial.ID,
		MetodoPagoID:         dto.MetodoPagoID,
		QrID:                 qr.ID,
	}

	if err := s.db.Create(&radicado).Error; err != nil {
		return nil, err
	}

	// ── 10. Preload completo ──
	if err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("DocumentoComercial.Area").
		Preload("DocumentoComercial.Moneda").
		Preload("DocumentoComercial.Detalles").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("Estado").
		Preload("MetodoPago").
		Preload("Qr").
		First(&radicado, radicado.ID).Error; err != nil {
		return nil, err
	}

	return &radicado, nil
}

func (s *Service) List() ([]db.DocumentoRadicado, error) {
	var radicados []db.DocumentoRadicado
	err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("Estado").
		Preload("MetodoPago").
		Preload("Qr").
		Order("fecha_radicacion DESC").
		Find(&radicados).Error
	return radicados, err
}

func (s *Service) GetByID(id uint) (*db.DocumentoRadicado, error) {
	var radicado db.DocumentoRadicado
	err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("DocumentoComercial.Area").
		Preload("DocumentoComercial.Moneda").
		Preload("DocumentoComercial.Detalles").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("UsuarioActual").
		Preload("PasoActual").
		Preload("Estado").
		Preload("MetodoPago").
		Preload("Qr").
		First(&radicado, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("el radicado no existe")
		}
		return nil, err
	}
	return &radicado, nil
}
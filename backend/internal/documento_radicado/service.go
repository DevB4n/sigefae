package documento_radicado

import (
	"errors"
	"fmt"
	"sort"
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
	var radicado db.DocumentoRadicado

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// ── 1. Validar que el documento comercial existe ──
		var docCom db.DocumentoComercial
		if err := tx.First(&docCom, dto.DocumentoComercialID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("el documento comercial no existe")
			}
			return err
		}

		// ── 2. Validar que NO esté ya radicado ──
		var existing db.DocumentoRadicado
		if err := tx.Where("documento_comercial_id = ?", dto.DocumentoComercialID).First(&existing).Error; err == nil {
			return errors.New("el documento ya fue radicado anteriormente")
		}

		// ── 3. Validar catálogos ──
		var tipoRad db.TipoRadicacion
		if err := tx.First(&tipoRad, dto.TipoRadicacionID).Error; err != nil {
			return errors.New("el tipo de radicación no existe")
		}

		var ruta db.Ruta
		if err := tx.First(&ruta, dto.RutaID).Error; err != nil {
			return errors.New("la ruta no existe")
		}

		var metodoPago db.MetodoPago
		if err := tx.First(&metodoPago, dto.MetodoPagoID).Error; err != nil {
			return errors.New("el método de pago no existe")
		}

		// ── 4. Obtener primer paso de la ruta ──
		var primerPaso db.PasoRuta
		var primerPasoID *uint
		responsableID := usuarioID
		estadoPosesion := "Libre"

		err := tx.Where("ruta_id = ?", dto.RutaID).Order("orden ASC, id ASC").First(&primerPaso).Error
		if err == nil {
			id := primerPaso.ID
			primerPasoID = &id
			if primerPaso.UsuarioID != 0 {
				responsableID = primerPaso.UsuarioID
			}
			estadoPosesion = "EnProceso"
		}

		// ── 5. Obtener estado inicial ──
		var estadoInicial db.EstadoDocumentoRadicado
		if err := tx.Order("id ASC").First(&estadoInicial).Error; err != nil {
			return errors.New("no hay estados de radicación configurados")
		}

		// ── 6. Generar número de radicado ──
		numeroRadicado := dto.NumeroRadicado
		if numeroRadicado == "" {
			year := time.Now().Year()
			var count int64
			tx.Model(&db.DocumentoRadicado{}).Where("YEAR(fecha_radicacion) = ?", year).Count(&count)
			numeroRadicado = fmt.Sprintf("RAD-%d-%05d", year, count+1)
		}

		var existingNum db.DocumentoRadicado
		if err := tx.Where("numero_radicado = ?", numeroRadicado).First(&existingNum).Error; err == nil {
			return errors.New("el número de radicado ya existe")
		}

		// ── 7. Crear QR ──
		qr := db.CodigoQr{
			Url:    fmt.Sprintf("https://sigefae.com/radicado/%s", numeroRadicado),
			Activo: true,
		}
		if err := tx.Create(&qr).Error; err != nil {
			return err
		}

		// ── 8. Crear radicado ──
		radicado = db.DocumentoRadicado{
			DocumentoComercialID: dto.DocumentoComercialID,
			TipoRadicacionID:     dto.TipoRadicacionID,
			RutaID:               dto.RutaID,
			NumeroRadicado:       numeroRadicado,
			FechaRadicacion:      time.Now(),
			UsuarioActualID:      responsableID,
			EstadoPosesion:       estadoPosesion,
			PasoActualID:         primerPasoID,
			EstadoID:             estadoInicial.ID,
			MetodoPagoID:         dto.MetodoPagoID,
			QrID:                 qr.ID,
		}

		if err := tx.Create(&radicado).Error; err != nil {
			return err
		}

		// ── 9. Generar tareas desde la ruta ──
		if err := generarTareasDesdeRuta(tx, &radicado, dto.RutaID, dto.DocumentoComercialID); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// ── 10. Preload completo fuera de la transacción ──
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
		Preload("PasoActual.Usuario").
		Preload("Estado").
		Preload("MetodoPago").
		Preload("Qr").
		Preload("Archivos").
		Preload("Archivos.Origen").
		First(&radicado, radicado.ID).Error; err != nil {
		return nil, err
	}

	return &radicado, nil
}

func (s *Service) List() ([]db.DocumentoRadicado, error) {
	var radicados []db.DocumentoRadicado
	if err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("MetodoPago").
		Preload("Estado").
		Preload("PasoActual").
		Preload("UsuarioActual").
		Preload("Archivos").
		Preload("Qr").
		Order("fecha_radicacion DESC").
		Find(&radicados).Error; err != nil {
		return nil, err
	}
	return radicados, nil
}

func (s *Service) GetByID(id uint) (*db.DocumentoRadicado, error) {
	var radicado db.DocumentoRadicado
	if err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("MetodoPago").
		Preload("Estado").
		Preload("PasoActual").
		Preload("UsuarioActual").
		Preload("Archivos").
		Preload("Qr").
		First(&radicado, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("documento radicado no encontrado")
		}
		return nil, err
	}
	return &radicado, nil
}

func (s *Service) Update(id uint, dto UpdateDTO) (*db.DocumentoRadicado, error) {
	var radicado db.DocumentoRadicado
	if err := s.db.First(&radicado, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("documento radicado no encontrado")
		}
		return nil, err
	}

	updates := map[string]interface{}{}
	if dto.TipoRadicacionID != 0 {
		updates["tipo_radicacion_id"] = dto.TipoRadicacionID
	}
	if dto.RutaID != 0 {
		updates["ruta_id"] = dto.RutaID
	}
	if dto.MetodoPagoID != 0 {
		updates["metodo_pago_id"] = dto.MetodoPagoID
	}
	if dto.NumeroRadicado != "" {
		updates["numero_radicado"] = dto.NumeroRadicado
	}
	if dto.UsuarioActualID != 0 {
		updates["usuario_actual_id"] = dto.UsuarioActualID
	}
	if dto.EstadoPosesion != "" {
		updates["estado_posesion"] = dto.EstadoPosesion
	}
	if dto.PasoActualID != 0 {
		updates["paso_actual_id"] = dto.PasoActualID
	}
	if dto.EstadoID != 0 {
		updates["estado_id"] = dto.EstadoID
	}

	if len(updates) > 0 {
		if err := s.db.Model(&radicado).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	if err := s.db.
		Preload("DocumentoComercial").
		Preload("DocumentoComercial.Proveedor").
		Preload("DocumentoComercial.Receptor").
		Preload("TipoRadicacion").
		Preload("Ruta").
		Preload("MetodoPago").
		Preload("Estado").
		Preload("PasoActual").
		Preload("UsuarioActual").
		Preload("Archivos").
		Preload("Qr").
		First(&radicado, id).Error; err != nil {
		return nil, err
	}
	return &radicado, nil
}

// ─────────────────────────────────────────────────────────────
// generarTareasDesdeRuta (FUNCIÓN PRIVADA DEL PACKAGE)
// ─────────────────────────────────────────────────────────────
func generarTareasDesdeRuta(tx *gorm.DB, radicado *db.DocumentoRadicado, rutaID uint, docComercialID uint) error {
	// 1. Traer pasos originales de la ruta
	var pasos []db.PasoRuta
	if err := tx.Where("ruta_id = ? AND activo = ?", rutaID, true).Order("orden asc").Find(&pasos).Error; err != nil {
		return err
	}

	// 2. Traer documento comercial para evaluar monto y área
	var docCom db.DocumentoComercial
	if err := tx.First(&docCom, docComercialID).Error; err != nil {
		return err
	}

	// 3. Evaluar reglas de monto que apliquen
	var reglas []db.ReglaMontoRuta
	tx.Where("activo = ? AND monto_minimo <= ?", true, docCom.Total).
		Where("(area_id IS NULL OR area_id = ?) AND (ruta_id IS NULL OR ruta_id = ?)", docCom.IDArea, rutaID).
		Preload("UsuarioAprobador").
		Find(&reglas)

	// 4. Construir pasos finales (originales + reglas)
	type pasoFinal struct {
		Orden      int
		Nombre     string
		UsuarioID  uint
		DesdeRegla bool
	}
	var finales []pasoFinal

	for _, p := range pasos {
		finales = append(finales, pasoFinal{Orden: p.Orden, Nombre: p.Nombre, UsuarioID: p.UsuarioID, DesdeRegla: false})
	}

	for _, r := range reglas {
		pf := pasoFinal{
			Nombre:     "Aprobación especial por monto",
			UsuarioID:  0,
			DesdeRegla: true,
		}
		if r.UsuarioAprobadorID != nil {
			pf.UsuarioID = *r.UsuarioAprobadorID
		} else if r.RolAprobadorID != nil {
			var usr db.Usuario
			if err := tx.Where("id_rol = ? AND activo = ?", *r.RolAprobadorID, true).First(&usr).Error; err == nil {
				pf.UsuarioID = usr.ID
			}
		}

		switch r.PosicionInsercion {
		case "PRIMERO":
			pf.Orden = -1000
			finales = append([]pasoFinal{pf}, finales...)
		case "ANTES_FINAL":
			pf.Orden = 9998
			finales = append(finales, pf)
		default: // ULTIMO
			pf.Orden = 9999
			finales = append(finales, pf)
		}
	}

	sort.Slice(finales, func(i, j int) bool {
		return finales[i].Orden < finales[j].Orden
	})

	// 5. Buscar estados base
	var estadoPendiente db.EstadoTarea
	tx.Where("nombre = ?", "Pendiente").First(&estadoPendiente)
	if estadoPendiente.ID == 0 {
		estadoPendiente.ID = 1
	}
	var estadoEnProceso db.EstadoTarea
	tx.Where("nombre = ?", "En Proceso").First(&estadoEnProceso)
	if estadoEnProceso.ID == 0 {
		estadoEnProceso.ID = 2
	}

	now := time.Now()

	// 6. Crear tareas
	for i, pf := range finales {
		uid := pf.UsuarioID
		if uid == 0 && i == 0 {
			uid = radicado.UsuarioActualID
		}

		tarea := db.Tarea{
			DocumentoRadicadoID: radicado.ID,
			UsuarioAsignadoID:   uid,
			EstadoID:            estadoPendiente.ID,
			Descripcion:         pf.Nombre,
			FechaAsignacion:     now,
		}
		if i == 0 {
			tarea.EstadoID = estadoEnProceso.ID
			tarea.FechaInicio = &now
			radicado.UsuarioActualID = uid
		}
		if err := tx.Create(&tarea).Error; err != nil {
			return err
		}
	}

	return tx.Save(radicado).Error
}
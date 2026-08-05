package documento_radicado

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
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

	// Validar normas de reparto
	if len(dto.NormasReparto) > 0 {
		var total float64
		for _, n := range dto.NormasReparto {
			total += n.Porcentaje
		}
		if total != 100 {
			return nil, errors.New("la suma de porcentajes de normas de reparto debe ser exactamente 100%")
		}
	}

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
			Url: fmt.Sprintf("http://localhost:5173/radicado/%s", numeroRadicado),
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

		// ── 10. Adjuntar archivos automáticamente desde el correo origen ──
		// Consultamos explícitamente correo_id para no depender del mapeo del struct
		var correoID *uint
		tx.Model(&db.DocumentoComercial{}).
			Select("correo_id").
			Where("id = ?", dto.DocumentoComercialID).
			Scan(&correoID)

		if correoID != nil && *correoID != 0 {
			if err := adjuntarArchivosDesdeCorreo(tx, &radicado, *correoID); err != nil {
				fmt.Printf("[RADICADO %s] warning: no se pudieron adjuntar archivos del correo: %v\n", radicado.NumeroRadicado, err)
			}
		} else {
			fmt.Printf("[RADICADO %s] info: documento comercial %d no tiene correo_id\n", radicado.NumeroRadicado, dto.DocumentoComercialID)
		}

		// Crear normas de reparto
		for _, n := range dto.NormasReparto {
			nr := db.RadicadoNormaReparto{
				DocumentoRadicadoID: radicado.ID,
				NormaRepartoID:      n.NormaRepartoID,
				Porcentaje:          n.Porcentaje,
			}
			if err := tx.Create(&nr).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// ── 11. Preload completo fuera de la transacción ──
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
		Preload("NormasReparto").
		Preload("NormasReparto.NormaReparto").
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
		Preload("NormasReparto").
		Preload("NormasReparto.NormaReparto").
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
		Preload("NormasReparto").
		Preload("NormasReparto.NormaReparto").
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
		Preload("NormasReparto").
		Preload("NormasReparto.NormaReparto").
		First(&radicado, id).Error; err != nil {
		return nil, err
	}
	return &radicado, nil
}

// ─────────────────────────────────────────────────────────────
// generarTareasDesdeRuta
// ─────────────────────────────────────────────────────────────
func generarTareasDesdeRuta(tx *gorm.DB, radicado *db.DocumentoRadicado, rutaID uint, docComercialID uint) error {
	// ── 1. Pasos base de la ruta ──
	var pasos []db.PasoRuta
	if err := tx.Where("ruta_id = ? AND activo = ?", rutaID, true).Order("orden asc, id asc").Find(&pasos).Error; err != nil {
		return err
	}
	if len(pasos) == 0 {
		return errors.New("la ruta no tiene pasos configurados")
	}

	// ── 2. Documento comercial (monto, área, moneda) ──
	var docCom db.DocumentoComercial
	if err := tx.First(&docCom, docComercialID).Error; err != nil {
		return err
	}

	// ── 3. Reglas de monto aplicables ──
	var reglas []db.ReglaMontoRuta
	query := tx.Where("activo = ? AND monto_minimo <= ?", true, docCom.Total).
		Where("(area_id IS NULL OR area_id = ?)", docCom.IDArea).
		Where("(ruta_id IS NULL OR ruta_id = ?)", rutaID).
		Preload("UsuarioAprobador").
		Preload("RolAprobador").
		Order("monto_minimo desc")

	if err := query.Find(&reglas).Error; err != nil {
		return err
	}

	// ── 4. Clasificar reglas por posición ──
	var alInicio, antesDelFinal, alFinal []db.ReglaMontoRuta
	for _, r := range reglas {
		switch r.PosicionInsercion {
		case "PRIMERO":
			alInicio = append(alInicio, r)
		case "ANTES_FINAL":
			antesDelFinal = append(antesDelFinal, r)
		default: // ULTIMO
			alFinal = append(alFinal, r)
		}
	}

	// ── 5. Helper para resolver usuario de una regla ──
	resolverUsuario := func(r db.ReglaMontoRuta) uint {
		if r.UsuarioAprobadorID != nil {
			return *r.UsuarioAprobadorID
		}
		if r.RolAprobadorID != nil {
			var usr db.Usuario
			if err := tx.Where("id_rol = ? AND activo = ?", *r.RolAprobadorID, true).First(&usr).Error; err == nil {
				return usr.ID
			}
		}
		return 0
	}

	type pasoFinal struct {
		Nombre    string
		UsuarioID uint
		EsRegla   bool
	}

	var flujo []pasoFinal

	// 5.1 PRIMERO
	for _, r := range alInicio {
		flujo = append(flujo, pasoFinal{
			Nombre:    fmt.Sprintf("Aprobación por monto ≥ %.2f", r.MontoMinimo),
			UsuarioID: resolverUsuario(r),
			EsRegla:   true,
		})
	}

	// 5.2 Pasos base (todos menos el último)
	for i := 0; i < len(pasos)-1; i++ {
		p := pasos[i]
		flujo = append(flujo, pasoFinal{
			Nombre:    p.Nombre,
			UsuarioID: p.UsuarioID,
			EsRegla:   false,
		})
	}

	// 5.3 ANTES_FINAL: justo antes del cierre del flujo base
	for _, r := range antesDelFinal {
		flujo = append(flujo, pasoFinal{
			Nombre:    fmt.Sprintf("Aprobación por monto ≥ %.2f", r.MontoMinimo),
			UsuarioID: resolverUsuario(r),
			EsRegla:   true,
		})
	}

	// 5.4 Último paso base (el cierre)
	ultimoPaso := pasos[len(pasos)-1]
	flujo = append(flujo, pasoFinal{
		Nombre:    ultimoPaso.Nombre,
		UsuarioID: ultimoPaso.UsuarioID,
		EsRegla:   false,
	})

	// 5.5 ULTIMO: después de todo
	for _, r := range alFinal {
		flujo = append(flujo, pasoFinal{
			Nombre:    fmt.Sprintf("Aprobación por monto ≥ %.2f", r.MontoMinimo),
			UsuarioID: resolverUsuario(r),
			EsRegla:   true,
		})
	}

	// ── 6. Estados base ──
	var estadoPendiente db.EstadoTarea
	if err := tx.Where("nombre = ?", "Pendiente").First(&estadoPendiente).Error; err != nil {
		return errors.New("no se encontró el estado 'Pendiente' para tareas")
	}
	var estadoEnProceso db.EstadoTarea
	if err := tx.Where("nombre = ?", "En Proceso").First(&estadoEnProceso).Error; err != nil {
		return errors.New("no se encontró el estado 'En Proceso' para tareas")
	}

	now := time.Now()

	// ── 7. Crear tareas en orden ──
	for i, pf := range flujo {
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

// ─────────────────────────────────────────────────────────────
// adjuntarArchivosDesdeCorreo
// ─────────────────────────────────────────────────────────────
func adjuntarArchivosDesdeCorreo(tx *gorm.DB, radicado *db.DocumentoRadicado, correoID uint) error {
	// 1. Cargar el correo para obtener id_mensaje
	var correo db.Correo
	if err := tx.First(&correo, correoID).Error; err != nil {
		return fmt.Errorf("correo origen no encontrado (id=%d): %w", correoID, err)
	}
	if correo.IDMensaje == "" {
		return errors.New("correo no tiene id_mensaje definido")
	}

	// 2. Buscar origen "Sistema"
	var origen db.ArchivoOrigen
	origenID := uint(0)
	if err := tx.Where("nombre = ? AND activo = ?", "Sistema", true).First(&origen).Error; err != nil {
		fmt.Printf("[RADICADO %d] WARNING: origen 'Sistema' no encontrado: %v\n", radicado.ID, err)
	} else {
		origenID = origen.ID
	}

	// 3. Directorio fuente
	srcDir := filepath.Join("storage", "mails", correo.IDMensaje)
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("no se pudo leer directorio origen %s: %w", srcDir, err)
	}

	// 4. Directorio destino
	dstDir := filepath.Join("storage", "radicados", fmt.Sprintf("%d", radicado.ID))
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio destino %s: %w", dstDir, err)
	}

		// 5. Copiar cada archivo (solo PDF)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		ext := strings.ToLower(filepath.Ext(name))

		// ← Solo archivos PDF
		if ext != ".pdf" {
			continue
		}

		srcPath := filepath.Join(srcDir, name)
		dstPath := filepath.Join(dstDir, name)

		if err := copyFile(srcPath, dstPath); err != nil {
			fmt.Printf("[RADICADO %d] ERROR copiando %s: %v\n", radicado.ID, name, err)
			continue
		}

		info, _ := os.Stat(dstPath)
		peso := int64(0)
		if info != nil {
			peso = info.Size()
		}

		archivo := db.Archivo{
			DocumentoRadicadoID: radicado.ID,
			Nombre:              name,
			Extension:           strings.TrimPrefix(ext, "."),
			Ruta:                dstPath,
			Peso:                peso,
			OrigenID:            origenID,
		}
		if err := tx.Create(&archivo).Error; err != nil {
			fmt.Printf("[RADICADO %d] ERROR guardando registro BD de %s: %v\n", radicado.ID, name, err)
			continue
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func (s *Service) GetByNumeroRadicado(numero string) (*db.DocumentoRadicado, error) {
    var radicado db.DocumentoRadicado
    if err := s.db.
        Preload("DocumentoComercial").
        Preload("DocumentoComercial.Proveedor").
        Preload("DocumentoComercial.Receptor").
        Preload("DocumentoComercial.Moneda").
        Preload("TipoRadicacion").
        Preload("Ruta").
        Preload("MetodoPago").
        Preload("Estado").
        Preload("PasoActual").
        Preload("UsuarioActual").
        Preload("Qr").
        Preload("NormasReparto").
        Preload("NormasReparto.NormaReparto").
        Where("numero_radicado = ?", numero).
        First(&radicado).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("documento no encontrado")
        }
        return nil, err
    }
    return &radicado, nil
}

func (s *Service) GetNormasReparto(radicadoID uint) ([]db.RadicadoNormaReparto, error) {
	var normas []db.RadicadoNormaReparto
	if err := s.db.Where("documento_radicado_id = ?", radicadoID).
		Preload("NormaReparto").
		Find(&normas).Error; err != nil {
		return nil, err
	}
	return normas, nil
}

func (s *Service) AsignarNormasReparto(radicadoID uint, dtos []NormaRepartoInputDTO) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var total float64
		for _, d := range dtos {
			total += d.Porcentaje
		}
		if total != 100 {
			return errors.New("la suma de porcentajes debe ser exactamente 100%")
		}

		if err := tx.Where("documento_radicado_id = ?", radicadoID).Delete(&db.RadicadoNormaReparto{}).Error; err != nil {
			return err
		}

		for _, d := range dtos {
			nr := db.RadicadoNormaReparto{
				DocumentoRadicadoID: radicadoID,
				NormaRepartoID:      d.NormaRepartoID,
				Porcentaje:          d.Porcentaje,
			}
			if err := tx.Create(&nr).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
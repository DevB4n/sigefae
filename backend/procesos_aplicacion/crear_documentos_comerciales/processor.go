package crear_documentos_comerciales

import (
	"errors"
	"log"
	"path/filepath"
	"time"

	"gorm.io/gorm"
	"sigefae/internal/db"
)

// ProcesarCarpetaCorreo escanea la carpeta de un correo buscando un .xml de DIAN y lo procesa.
func ProcesarCarpetaCorreo(database *gorm.DB, correoID uint, folderPath string) error {
	matches, err := filepath.Glob(filepath.Join(folderPath, "*.xml"))
	if err != nil || len(matches) == 0 {
		// No hay XML, puede ser otro tipo de adjunto. No es un error crítico.
		return nil
	}

	for _, xmlFile := range matches {
		invoice, err := ParseXMLFile(xmlFile)
		if err != nil {
			log.Printf("Error parseando XML %s: %v\n", xmlFile, err)
			continue
		}

		err = ProcesarFacturaXML(database, invoice, correoID)
		if err != nil {
			log.Printf("Error registrando factura desde XML %s: %v\n", xmlFile, err)
			continue
		}

		// Asegurar que el estado 2 (Procesado) exista
		var estado db.EstadoCorreo
		if err := database.Where("id = ?", 2).FirstOrCreate(&estado, db.EstadoCorreo{ID: 2, Nombre: "Procesado", Activo: true}).Error; err != nil {
			log.Printf("Error creando estado Procesado para el correo %d: %v\n", correoID, err)
		}

		database.Model(&db.Correo{}).Where("id = ?", correoID).Update("id_estado", 2) // Asumiendo estado 2 = Procesado
		break
	}
	return nil
}

// ProcesarFacturaXML lee la estructura XML y nutre la base de datos (Proveedor, Receptor, Documento)
func ProcesarFacturaXML(database *gorm.DB, invoice *Invoice, correoID uint) error {
	return database.Transaction(func(tx *gorm.DB) error {

		// 0. Garantizar la existencia de catálogos base requeridos para foráneas
		var tipoDoc db.TipoDocumento
		if err := tx.Where("nombre = ?", "NIT").FirstOrCreate(&tipoDoc, db.TipoDocumento{Nombre: "NIT", Activo: true}).Error; err != nil {
			return err
		}

		var categoria db.CategoriaProveedor
		if err := tx.Where("nombre = ?", "General").FirstOrCreate(&categoria, db.CategoriaProveedor{Nombre: "General", Descripcion: "Generado automáticamente", Activo: true}).Error; err != nil {
			return err
		}

		var tipoPersona db.TipoPersona
		if err := tx.Where("nombre = ?", "Jurídica").FirstOrCreate(&tipoPersona, db.TipoPersona{Nombre: "Jurídica", Activo: true}).Error; err != nil {
			return err
		}

		var actividad db.ActividadEconomica
		if err := tx.Where("codigo = ?", "0000").FirstOrCreate(&actividad, db.ActividadEconomica{Nombre: "General", Codigo: "0000", Activo: true}).Error; err != nil {
			return err
		}

		var pais db.Pais
		if err := tx.Where("codigo = ?", "CO").FirstOrCreate(&pais, db.Pais{Nombre: "Colombia", Codigo: "CO", Activo: true}).Error; err != nil {
			return err
		}

		var moneda db.Moneda
		if err := tx.Where("codigo = ?", "COP").FirstOrCreate(&moneda, db.Moneda{Nombre: "Peso Colombiano", Codigo: "COP", Activo: true}).Error; err != nil {
			return err
		}

		var area db.Area
		if err := tx.Where("nombre = ?", "General").FirstOrCreate(&area, db.Area{Nombre: "General", Activo: true}).Error; err != nil {
			return err
		}

		// 1. Gestionar Receptor (Quien recibe la factura - La Empresa)
		nitReceptor := invoice.AccountingCustomerParty.Party.PartyTaxScheme.CompanyID
		var receptor db.Receptor
		if err := tx.Where("numero_documento = ?", nitReceptor).First(&receptor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Crear receptor si no existe
				receptor = db.Receptor{
					Nombre:          invoice.AccountingCustomerParty.Party.PartyTaxScheme.RegistrationName,
					NumeroDocumento: nitReceptor,
					TipoDocumentoID: tipoDoc.ID,
					Activo:          true,
				}
				if err := tx.Create(&receptor).Error; err != nil {
					return err
				}
			} else {
				return err
			}
		}

		// 2. Gestionar Geografía (País, Depto, Municipio)
		supplierParty := invoice.AccountingSupplierParty.Party
		addr := supplierParty.PhysicalLocation.Address

		// 3. Gestionar Proveedor
		nitProveedor := supplierParty.PartyTaxScheme.CompanyID
		var proveedor db.Proveedor
		if err := tx.Where("numero_documento = ?", nitProveedor).First(&proveedor).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Buscar o Crear Geografía para el Proveedor
				var depto db.Departamento
				tx.Where("nombre = ?", addr.CountrySubentity).FirstOrCreate(&depto, db.Departamento{Nombre: addr.CountrySubentity, PaisID: pais.ID})

				var municipio db.Municipio
				tx.Where("nombre = ?", addr.CityName).FirstOrCreate(&municipio, db.Municipio{Nombre: addr.CityName, DepartamentoID: depto.ID})

				var direccion db.Direccion
				tx.Where("linea_1 = ? AND id_municipio = ?", addr.AddressLine.Line, municipio.ID).FirstOrCreate(&direccion, db.Direccion{
					Nombre:      "Principal",
					Linea1:      addr.AddressLine.Line,
					IDMunicipio: municipio.ID,
				})

				// Crear Proveedor
				proveedor = db.Proveedor{
					NumeroDocumento:      nitProveedor,
					TipoDocumentoID:      tipoDoc.ID,
					RazonSocial:          supplierParty.PartyTaxScheme.RegistrationName,
					NombreComercial:      supplierParty.PartyName.Name,
					DireccionID:          direccion.ID,
					CategoriaID:          categoria.ID,
					TipoPersonaID:        tipoPersona.ID,
					ActividadEconomicaID: actividad.ID,
					Activo:               true,
				}
				if err := tx.Create(&proveedor).Error; err != nil {
					return err
				}

				// Crear Contacto si existe
				if supplierParty.Contact.Name != "" || supplierParty.Contact.ElectronicMail != "" {
					contacto := db.Contacto{
						ProveedorID: proveedor.ID,
						Nombre:      supplierParty.Contact.Name,
						Correo:      supplierParty.Contact.ElectronicMail,
						Telefono:    supplierParty.Contact.Telephone,
						Activo:      true,
					}
					tx.Create(&contacto)
				}

				// Responsabilidad Fiscal
				if supplierParty.PartyTaxScheme.TaxLevelCode != "" {
					resp := db.ResponsabilidadFiscal{
						IDProveedor: proveedor.ID,
						Codigo:      supplierParty.PartyTaxScheme.TaxLevelCode,
						Activo:      true,
					}
					tx.Create(&resp)
				}
			} else {
				return err
			}
		}

		// 4. Crear Documento Comercial

		// Validar si la factura ya existe (CUFE)
		var docExistente db.DocumentoComercial
		if err := tx.Where("cufe = ?", invoice.UUID).First(&docExistente).Error; err == nil {
			// Si ya existe por CUFE, la saltamos.
			return nil
		}

		issueDate, _ := time.Parse("2006-01-02", invoice.IssueDate)
		var dueDate *time.Time
		if invoice.DueDate != "" {
			d, err := time.Parse("2006-01-02", invoice.DueDate)
			if err == nil {
				dueDate = &d
			}
		}

		// Subtotal: preferir LineExtensionAmount (suma de líneas).
		// Si viene en 0, calcular desde los detalles como fallback.
		subtotal := invoice.LegalMonetaryTotal.LineExtensionAmount
		if subtotal == 0 && len(invoice.InvoiceLines) > 0 {
			for _, line := range invoice.InvoiceLines {
				subtotal += line.LineExtensionAmount
			}
		}

		// IVA: buscar explícitamente el tax scheme "01" (IVA) para no confundir
		// con retenciones u otros impuestos que puedan aparecer en TaxTotal.
		var iva float64
		for _, tax := range invoice.TaxTotal {
			// Nota: si en tu struct de parser agregas el TaxScheme ID, puedes filtrar
			// tax.TaxScheme.ID == "01". Mientras tanto sumamos TaxAmount.
			iva += tax.TaxAmount
		}

		doc := db.DocumentoComercial{
			Tipo:             "FACTURA_ELECTRONICA",
			NumeroDocumento:  invoice.ID,
			IDProveedor:      proveedor.ID,
			IDReceptor:       receptor.ID,
			IDArea:           area.ID,
			FechaDocumento:   issueDate,
			FechaVencimiento: dueDate,
			MonedaID:         moneda.ID,
			Subtotal:         subtotal,
			Iva:              iva,
			Total:            invoice.LegalMonetaryTotal.PayableAmount,
			Cufe:             &invoice.UUID,
			CorreoID:         &correoID,
			Activo:           true,
		}

		if err := tx.Create(&doc).Error; err != nil {
			return err
		}

		// 5. Insertar Detalles de Factura
		for _, line := range invoice.InvoiceLines {
			// Calcular valor unitario neto desde LineExtensionAmount / Cantidad
			// porque PriceAmount puede venir con IVA incluido en algunos proveedores
			valorUnit := line.Price.PriceAmount
			if line.InvoicedQuantity > 0 {
				valorUnit = line.LineExtensionAmount / line.InvoicedQuantity
			}

			detalle := db.DetalleDocumentoComercial{
				DocumentoComercialID: doc.ID,
				Descripcion:          line.Item.Description,
				Cantidad:             line.InvoicedQuantity,
				ValorUnit:            valorUnit,
				Total:                line.LineExtensionAmount,
				Activo:               true,
			}
			if err := tx.Create(&detalle).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

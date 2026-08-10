package db

import (
	"time"

	"gorm.io/gorm"
)

func strPtr(s string) *string { return &s }

func SeedNormasReparto(db *gorm.DB) error {
	var count int64
	db.Model(&NormaReparto{}).Count(&count)
	if count > 0 {
		return nil
	}

	normas := []NormaReparto{
		// BUCARAMANGA - ADMON
		{Codigo: "BU101", Nombre: "GERENCIA", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU102", Nombre: "DIRECCION ADMINIST Y F/CIERA", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU103", Nombre: "GESTION HUMANA", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU104", Nombre: "DIRECCION CONTABLE E IMPUESTOS", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU106", Nombre: "SERVICIOS GENERALES", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU108", Nombre: "SISTEMAS", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU110", Nombre: "GESTION AMBIENTAL", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU112", Nombre: "DIR ADMON Y COMPRAS", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU120", Nombre: "GASTOS GENERALES", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU122", Nombre: "SALUD OCUPACIONAL", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU123", Nombre: "JUNTA DIRECTIVA", Sucursal: "BUCARAMANGA", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},

		// BUCARAMANGA - VENTAS
		{Codigo: "BU701", Nombre: "SUBGERENCIA GENERAL", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU702", Nombre: "DIRECCION COMERCIAL", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU703", Nombre: "MERCADEO BGA", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU704", Nombre: "CARTERA", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU706", Nombre: "LOGISTICA Y TRANSPORTE", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU710", Nombre: "GESTION AMBIENTAL", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU716", Nombre: "SERVICIO AL CLIENTE", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU720", Nombre: "GASTOS GENERALES", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU722", Nombre: "SALUD OCUPACIONAL", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU724", Nombre: "ANALISTA DE COMPRAS Y LOGISTIC", Sucursal: "BUCARAMANGA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},

		// BUCARAMANGA - PRODUCCION
		{Codigo: "BU901", Nombre: "DIRECCION DE PRODUCCION", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU902", Nombre: "LABORATORIO", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU903", Nombre: "RECIBO DE MATERIA PRIMA", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU904", Nombre: "ACONDICIONAMIENTO DE TRIGO", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU905", Nombre: "PROCESO DE MOLIENDA", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU906", Nombre: "PROCESO DE EMPAQUE", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU907", Nombre: "PROCESO ROBINSON", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU908", Nombre: "MTTO ELECTRICO Y MECANICO", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU910", Nombre: "GESTION AMBIENTAL", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU911", Nombre: "INVESTIGACION Y DESARROLLO", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU912", Nombre: "LOGISTICA DE IMPORTACIONES", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU920", Nombre: "GASTOS GENERALES", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "BU922", Nombre: "SALUD OCUPACIONAL", Sucursal: "BUCARAMANGA", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},

		// MALAMBO - ADMON
		{Codigo: "MB101", Nombre: "GERENCIA", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB102", Nombre: "DIRECCION ADMINIST Y F/CIERA", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB103", Nombre: "GESTION HUMANA", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB104", Nombre: "DIRECCION CONTABLE E IMPUESTOS", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB106", Nombre: "SERVICIOS GENERALES", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB108", Nombre: "SISTEMAS", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB109", Nombre: "SISTEMA DE GESTION DE CALIDAD", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB110", Nombre: "GESTION AMBIENTAL", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB112", Nombre: "DIR ADMON Y COMPRAS", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB120", Nombre: "GASTOS GENERALES", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB122", Nombre: "SALUD OCUPACIONAL", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB123", Nombre: "JUNTA DIRECTIVA", Sucursal: "MALAMBO", Departamento: "ADMON", Activo: true, CreatedAt: time.Now()},

		// MALAMBO - VENTAS
		{Codigo: "MB701", Nombre: "SUBGERENCIAL GENERAL", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB702", Nombre: "DIRECCION COMERCIAL", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB703", Nombre: "MERCADEO MB", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB704", Nombre: "CARTERA", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB706", Nombre: "LOGISTICA", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB707", Nombre: "SERVICIO AL CLIENTE", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB710", Nombre: "GESTION AMBIENTAL", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB716", Nombre: "SERVICIO AL CLIENTE", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB720", Nombre: "GASTOS GENERALES", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB722", Nombre: "SALUD OCUPACIONAL", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB724", Nombre: "ANALISTA DE COMPRAS Y LOGISTIC", Sucursal: "MALAMBO", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},

		// MALAMBO - PRODUCCION
		{Codigo: "MB901", Nombre: "DIRECCION PRODUCCION", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB902", Nombre: "LABORATORIO", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB903", Nombre: "RECIBO DE MATERIA PRIMA", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB904", Nombre: "ACONDICIONAMIENTO DE TRIGO", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB906", Nombre: "PROCESO DE EMPAQUE", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB907", Nombre: "PROCESO ROBINSON", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB908", Nombre: "MTTO ELECTRICO Y MECANICO", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB909", Nombre: "SISTEMA DE GESTION DE CALIDAD", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB910", Nombre: "GESTION AMBIENTAL", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB911", Nombre: "INVESTIGACION Y DESARROLLO", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB912", Nombre: "LOGISTICA DE IMPORTACIONES", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB920", Nombre: "GASTOS GENERALES", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},
		{Codigo: "MB922", Nombre: "SALUD OCUPACIONAL", Sucursal: "MALAMBO", Departamento: "PRODUCCION", Activo: true, CreatedAt: time.Now()},

		// CUCUTA - VENTAS
		{Codigo: "CU702", Nombre: "DIRECCION COMERCIAL", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU703", Nombre: "MERCADEO", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU704", Nombre: "CARTERA", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU706", Nombre: "LOGISTICA Y TRANSPORTE", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU710", Nombre: "AGENCIA CUCUTA", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU722", Nombre: "SALUD OCUPACIONAL", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CU724", Nombre: "ANALISTA DE COMPRAS Y LOGISTIC", Sucursal: "CUCUTA", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},

		// CB - VENTAS
		{Codigo: "CB702", Nombre: "DIRECCION COMERCIAL", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CB704", Nombre: "CARTERA", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CB705", Nombre: "AGENCIA CENTROABASTOS", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CB706", Nombre: "LOGISTICA Y TRANSPORTE", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CB722", Nombre: "SALUD OCUPACIONAL", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},
		{Codigo: "CB724", Nombre: "ANALISTA DE COMPRAS Y LOGISTIC", Sucursal: "CB", Departamento: "VENTAS", Activo: true, CreatedAt: time.Now()},

		// INDICADORES IMPUESTOS - ADMON
		{Codigo: "IVATRC33", Nombre: "IVA TRANS RC ADMON SERV TARIFA 19%", Sucursal: "GENERAL", Departamento: "ADMON", Tipo: strPtr("Servicio"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC32", Nombre: "IVA TRANS RC ADMON SERV TARIFA 5%", Sucursal: "GENERAL", Departamento: "ADMON", Tipo: strPtr("Servicio"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC29", Nombre: "IVA TRANS RC ADMON TARIFA 19%", Sucursal: "GENERAL", Departamento: "ADMON", Tipo: strPtr("Compra"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC23", Nombre: "IVA TRANS RC ADMON TARIFA 5%", Sucursal: "GENERAL", Departamento: "ADMON", Tipo: strPtr("Compra"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},

		// INDICADORES IMPUESTOS - VENTAS
		{Codigo: "IVATRC35", Nombre: "IVA TRANS RC VENTAS SERV TARIFA 19%", Sucursal: "GENERAL", Departamento: "VENTAS", Tipo: strPtr("Servicio"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC34", Nombre: "IVA TRANS RC VENTAS SERV TARIFA 5%", Sucursal: "GENERAL", Departamento: "VENTAS", Tipo: strPtr("Servicio"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC30", Nombre: "IVA TRANS RC VENTAS TARIFA 19%", Sucursal: "GENERAL", Departamento: "VENTAS", Tipo: strPtr("Compra"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC24", Nombre: "IVA TRANS RC VENTAS 5%", Sucursal: "GENERAL", Departamento: "VENTAS", Tipo: strPtr("Compra"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},

		// INDICADORES IMPUESTOS - PRODUCCION
		{Codigo: "IVATRC37", Nombre: "IVA TRANS RC PRODUCCION SERV TARIFA 19%", Sucursal: "GENERAL", Departamento: "PRODUCCION", Tipo: strPtr("Servicio"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC36", Nombre: "IVA TRANS RC PRODUCCION SERV TARIFA 5%", Sucursal: "GENERAL", Departamento: "PRODUCCION", Tipo: strPtr("Servicio"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC31", Nombre: "IVA TRANS RC PRODUCCION TARIFA 19%", Sucursal: "GENERAL", Departamento: "PRODUCCION", Tipo: strPtr("Compra"), TarifaIva: strPtr("19%"), Activo: true, CreatedAt: time.Now()},
		{Codigo: "IVATRC25", Nombre: "IVA TRANS RC PRODUCCION 5%", Sucursal: "GENERAL", Departamento: "PRODUCCION", Tipo: strPtr("Compra"), TarifaIva: strPtr("5%"), Activo: true, CreatedAt: time.Now()},
	}

	for _, n := range normas {
		if err := db.Create(&n).Error; err != nil {
			return err
		}
	}
	return nil
}

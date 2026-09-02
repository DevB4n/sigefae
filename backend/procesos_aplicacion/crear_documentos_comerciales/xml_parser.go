package crear_documentos_comerciales

import (
	"encoding/xml"
	"errors"
	"os"
	"strings"
)

// Invoice representa la estructura principal de una Factura Electrónica UBL 2.1 (DIAN).
type Invoice struct {
	XMLName xml.Name `xml:"Invoice"`
	ID      string   `xml:"ID"`
	UUID    string   `xml:"UUID"` // CUFE

	IssueDate string `xml:"IssueDate"`
	DueDate   string `xml:"DueDate"`

	PaymentMeans []struct {
		ID               string `xml:"ID"`
		PaymentMeansCode string `xml:"PaymentMeansCode"`
		PaymentDueDate   string `xml:"PaymentDueDate"`
	} `xml:"PaymentMeans"`

	// Emisor / Proveedor
	AccountingSupplierParty struct {
		Party struct {
			PartyName struct {
				Name string `xml:"Name"`
			} `xml:"PartyName"`
			PhysicalLocation struct {
				Address struct {
					CityName         string `xml:"CityName"`
					CountrySubentity string `xml:"CountrySubentity"` // Departamento
					AddressLine      struct {
						Line string `xml:"Line"`
					} `xml:"AddressLine"`
				} `xml:"Address"`
			} `xml:"PhysicalLocation"`
			PartyTaxScheme struct {
				RegistrationName string `xml:"RegistrationName"`
				CompanyID        string `xml:"CompanyID"`    // NIT
				TaxLevelCode     string `xml:"TaxLevelCode"` // Responsabilidad Fiscal
			} `xml:"PartyTaxScheme"`
			Contact struct {
				Name           string `xml:"Name"`
				Telephone      string `xml:"Telephone"`
				ElectronicMail string `xml:"ElectronicMail"`
			} `xml:"Contact"`
		} `xml:"Party"`
	} `xml:"AccountingSupplierParty"`

	// Receptor / Cliente
	AccountingCustomerParty struct {
		Party struct {
			PartyTaxScheme struct {
				RegistrationName string `xml:"RegistrationName"`
				CompanyID        string `xml:"CompanyID"` // NIT
			} `xml:"PartyTaxScheme"`
			PhysicalLocation struct {
				Address struct {
					CityName    string `xml:"CityName"`
					AddressLine struct {
						Line string `xml:"Line"`
					} `xml:"AddressLine"`
				} `xml:"Address"`
			} `xml:"PhysicalLocation"`
		} `xml:"Party"`
	} `xml:"AccountingCustomerParty"`

	// Totales
	LegalMonetaryTotal struct {
		LineExtensionAmount float64 `xml:"LineExtensionAmount"`
		TaxExclusiveAmount  float64 `xml:"TaxExclusiveAmount"`
		TaxInclusiveAmount  float64 `xml:"TaxInclusiveAmount"`
		PayableAmount       float64 `xml:"PayableAmount"`
	} `xml:"LegalMonetaryTotal"`

	TaxTotal []struct {
		TaxAmount float64 `xml:"TaxAmount"`
	} `xml:"TaxTotal"`

	// Detalle de los ítems facturados
	InvoiceLines []InvoiceLine `xml:"InvoiceLine"`
}

type InvoiceLine struct {
	ID                  string  `xml:"ID"`
	InvoicedQuantity    float64 `xml:"InvoicedQuantity"`
	LineExtensionAmount float64 `xml:"LineExtensionAmount"`
	Item                struct {
		Description string `xml:"Description"`
	} `xml:"Item"`
	Price struct {
		PriceAmount float64 `xml:"PriceAmount"`
	} `xml:"Price"`
}

// AttachedDocument representa el envoltorio estándar de la DIAN para facturas electrónicas.
type AttachedDocument struct {
	XMLName    xml.Name `xml:"AttachedDocument"`
	Attachment struct {
		ExternalReference struct {
			Description string `xml:"Description"`
		} `xml:"ExternalReference"`
	} `xml:"Attachment"`
	ParentDocumentLineReference struct {
		DocumentReference struct {
			Attachment struct {
				ExternalReference struct {
					Description string `xml:"Description"`
				} `xml:"ExternalReference"`
			} `xml:"Attachment"`
		} `xml:"DocumentReference"`
	} `xml:"ParentDocumentLineReference"`
}

// ParseXMLFile lee y parsea el archivo XML de la DIAN a la estructura Invoice
func ParseXMLFile(filepath string) (*Invoice, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// 1. Intentar decodificar como AttachedDocument (UBL Colombiano)
	var attached AttachedDocument
	if err := xml.Unmarshal(data, &attached); err == nil {
		// La factura puede estar en Attachment o en ParentDocumentLineReference
		cdataLists := []string{
			attached.Attachment.ExternalReference.Description,
			attached.ParentDocumentLineReference.DocumentReference.Attachment.ExternalReference.Description,
		}

		for _, cdataStr := range cdataLists {
			if cdataStr != "" {
				var invoice Invoice
				if err := xml.Unmarshal([]byte(cdataStr), &invoice); err == nil && invoice.ID != "" {
					return &invoice, nil
				}
			}
		}
	}

	// 2. Fallback: Intentar decodificar directamente como Invoice (XML crudo sin envoltorio)
	var invoice Invoice
	if err := xml.Unmarshal(data, &invoice); err != nil {
		if strings.Contains(err.Error(), "expected element type <Invoice>") {
			return nil, errors.New("el documento XML no es una Factura DIAN válida o no contiene una factura en su interior")
		}
		return nil, err
	}

	return &invoice, nil
}

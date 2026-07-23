package registro_aprobacion

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint      `json:"id"`
	DocumentoRadicadoID uint      `json:"documento_radicado_id"`
	ResponsableID       uint      `json:"responsable_id"`
	RolID               uint      `json:"rol_id"`
	Estado              string    `json:"estado"`
	Observacion         string    `json:"observacion"`
	Fecha               time.Time `json:"fecha"`
}

func toResponse(registro db.RegistroAprobacion) Response {

	return Response{
		ID:                  registro.ID,
		DocumentoRadicadoID: registro.DocumentoRadicadoID,
		ResponsableID:       registro.ResponsableID,
		RolID:               registro.RolID,
		Estado:              registro.Estado,
		Observacion:         registro.Observacion,
		Fecha:               registro.Fecha,
	}
}
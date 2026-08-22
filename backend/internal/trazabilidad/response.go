package trazabilidad

import (
	"time"

	"sigefae/internal/db"
)

type Response struct {
	ID                  uint      `json:"id"`
	DocumentoRadicadoID uint      `json:"documento_radicado_id"`
	UsuarioID           uint      `json:"usuario_id"`
	UsuarioNombre       string    `json:"usuario_nombre"`
	Accion              string    `json:"accion"`
	Descripcion         string    `json:"descripcion"`
	Fecha               time.Time `json:"fecha"`
}

func toResponse(trazabilidad db.Trazabilidad) Response {

	userName := ""
	if trazabilidad.Usuario != nil {
		userName = trazabilidad.Usuario.Nombre
	}

	return Response{
		ID:                  trazabilidad.ID,
		DocumentoRadicadoID: trazabilidad.DocumentoRadicadoID,
		UsuarioID:           trazabilidad.UsuarioID,
		UsuarioNombre:       userName,
		Accion:              trazabilidad.Accion,
		Descripcion:         trazabilidad.Descripcion,
		Fecha:               trazabilidad.Fecha,
	}
}
type RadicadoConTrazabilidadResponse struct {
	ID              uint                  `json:"id"`
	NumeroRadicado  string                `json:"numero_radicado"`
	FechaRadicacion time.Time             `json:"fecha_radicacion"`
	EstadoPosesion  string                `json:"estado_posesion"`
	Proveedor       ProveedorMiniResponse `json:"proveedor"`
	Area            AreaMiniResponse      `json:"area"`
	Trazabilidad    []Response            `json:"trazabilidad"`
}

type ProveedorMiniResponse struct {
	ID              uint   `json:"id"`
	RazonSocial     string `json:"razon_social"`
	NumeroDocumento string `json:"numero_documento"`
}

type AreaMiniResponse struct {
	ID     uint   `json:"id"`
	Nombre string `json:"nombre"`
}
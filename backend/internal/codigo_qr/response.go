package codigo_qr

import "sigefae/internal/db"

type Response struct {
	ID     uint   `json:"id"`
	Url    string `json:"url"`
	Activo bool   `json:"activo"`
}

func toResponse(qr db.CodigoQr) Response {

	return Response{
		ID:     qr.ID,
		Url:    qr.Url,
		Activo: qr.Activo,
	}
}

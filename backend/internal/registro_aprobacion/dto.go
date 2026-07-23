package registro_aprobacion

type CreateDTO struct {
	DocumentoRadicadoID uint   `json:"documento_radicado_id" binding:"required"`
	ResponsableID       uint   `json:"responsable_id" binding:"required"`
	RolID               uint   `json:"rol_id" binding:"required"`
	Estado              string `json:"estado" binding:"required"`
	Observacion         string `json:"observacion"`
}

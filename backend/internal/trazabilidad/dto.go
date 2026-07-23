package trazabilidad

type CreateDTO struct {
	DocumentoRadicadoID uint   `json:"documento_radicado_id" binding:"required"`
	UsuarioID           uint   `json:"usuario_id" binding:"required"`
	Accion              string `json:"accion" binding:"required"`
	Descripcion         string `json:"descripcion"`
}

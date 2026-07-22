package comentario

type CreateDTO struct {
	DocumentoRadicadoID uint   `json:"documento_radicado_id" binding:"required"`
	UsuarioID           uint   `json:"usuario_id" binding:"required"`
	Descripcion         string `json:"descripcion" binding:"required"`
}
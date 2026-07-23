package archivo

type CreateDTO struct {
	DocumentoRadicadoID uint   `json:"documento_radicado_id" binding:"required"`
	Nombre              string `json:"nombre" binding:"required"`
	Extension           string `json:"extension" binding:"required"`
	Peso                int64  `json:"peso" binding:"required"`
	Ruta                string `json:"ruta" binding:"required"`
	OrigenID            uint   `json:"origen_id" binding:"required"`
}

type UpdateDTO struct {
	Nombre    string `json:"nombre" binding:"required"`
	Extension string `json:"extension" binding:"required"`
	Peso      int64  `json:"peso" binding:"required"`
	Ruta      string `json:"ruta" binding:"required"`
	OrigenID  uint   `json:"origen_id" binding:"required"`
}

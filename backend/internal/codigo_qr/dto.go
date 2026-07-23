package codigo_qr

type CreateDTO struct {
	Url string `json:"url" binding:"required"`
}

type UpdateDTO struct {
	Url string `json:"url" binding:"required"`
}

type UpdateStatusDTO struct {
	Activo bool `json:"activo"`
}

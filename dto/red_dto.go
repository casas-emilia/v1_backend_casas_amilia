package dto

type CrearRedRequest struct {
	RedSocial string `json:"red_social" binding:"required"`
	Link      string `json:"link"`
	EmpresaID uint   `json:"empresa_id"`
}

type ActualizarRedRequest struct {
	RedSocial string `json:"red_social" binding:"required"`
	Link      string `json:"link"`
}

type RedResponse struct {
	ID        uint   `json:"id"`
	RedSocial string `json:"red_social" binding:"required"`
	Link      string `json:"link"`
	EmpresaID uint   `json:"empresa_id"`
}

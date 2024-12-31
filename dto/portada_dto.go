package dto

import "time"

type CrearPortadaRequest struct {
	NombrePortada string `form:"nombre_portada" binding:"required"` // Para multipart/form-data
	//Image         string `json:"image"`
	EmpresaID uint `form:"empresa_id"` // Si envías este campo en el formulario
}

type ActualizarPortadaRequest struct {
	NombrePortada string `form:"nombre_portada"`
	//Image         string `form:"image"`
	EmpresaID uint `form:"empresa_id"`
}

type PortadaResponse struct {
	ID            uint      `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	NombrePortada string    `json:"nombre_portada"`
	Image         string    `json:"image"`
	EmpresaID     uint      `json:"empresa_id"`
}

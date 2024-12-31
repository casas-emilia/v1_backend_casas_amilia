package dto

import "time"

type CrearContactoRequest struct {
	EmailLaboral     string `json:"email_laboral" binding:"required,email"`
	CelularLaboral   string `json:"celular_laboral" binding:"min=6"` // minimo seis números
	DireccionLaboral string `json:"direccion_laboral"`
	//UsuarioID        uint   `json:"usuario_id" binding:"required"`
}

type ActualizarContactoRequest struct {
	EmailLaboral     string `json:"email_laboral" binding:"required,email"`
	CelularLaboral   string `json:"celular_laboral" binding:"min=6"`
	DireccionLaboral string `json:"direccion_laboral"`
}

type ContactoResponse struct {
	ID               uint      `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	EmailLaboral     string    `json:"email_laboral" binding:"required,email"`
	CelularLaboral   string    `json:"celular_laboral" binding:"min=8"`
	DireccionLaboral string    `json:"direccion_laboral"`
	UsuarioID        uint      `json:"usuario_id" binding:"required"`
}

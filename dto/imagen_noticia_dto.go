package dto

import "time"

type CrearImagen_noticiaRequest struct {
	//Image string `json:"image" binding:"required"`
}

type ActualizarImagen_noticiaRequest struct {
	//Image string `json:"image" binding:"required"`
}

type Imagen_noticiaResponse struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Image     string    `json:"image"`
	NoticiaID uint      `json:"noticia_id"`
}

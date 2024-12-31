package dto

import "time"

type CrearCaracteristicaRequest struct {
	Clave string `json:"clave" binding:"required"`
	Valor string `json:"valor" binding:"required"`
	//PrefabricadaID uint   `json:"prefabricada_id" binding:"required"`
}

type ActualizarCaracteristicaRequest struct {
	Clave string `json:"clave" binding:"required"`
	Valor string `json:"valor" binding:"required"`
}

type CaracteristicaResponse struct {
	ID             uint      `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Clave          string    `json:"clave" binding:"required"`
	Valor          string    `json:"valor" binding:"required"`
	PrefabricadaID uint      `json:"prefabricada_id" binding:"required"`
}

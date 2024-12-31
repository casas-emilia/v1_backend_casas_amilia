package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminObtenerServicios(c *gin.Context) {
	// Simula datos del panel de administración
	response := gin.H{
		"stats": gin.H{
			"total_usuarios": 120,
			"ventas_mes":     5000,
			"nuevas_ordenes": 25,
		},
		"notificaciones": []string{
			"Tienes 5 usuarios nuevos registrados.",
			"Recuerda revisar las órdenes pendientes.",
		},
	}

	// Responde con JSON
	c.JSON(http.StatusOK, response)
}

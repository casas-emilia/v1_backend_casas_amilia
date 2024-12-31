package controllers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleError(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort() // Asegura que no se ejecute más lógica después.
}

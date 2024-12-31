package main

import (
	"log"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/routers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Cargar variables de entorno desde el archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No se pudo cargar el archivo .env, asegurarse de que las variables de entorno estén configuradas")
	}

	// Conectar a la base de datos
	configs.ConnectToDB()
}

func main() {
	// Configurar y correr el servidor
	router := routers.SetupRouter() // Llamar a la función que configura las rutas

	// Add a health check endpoint
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	router.Run(":8080")
}

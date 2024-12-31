package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función auxiliar para manejar errores
//
//	func handleErrorEstilo(c *gin.Context, err error, statusCode int, message string) {
//		if err != nil {
//			c.JSON(statusCode, gin.H{"error": message})
//		}
//	}
func handleErrorEstilo(c *gin.Context, err error, statusCode int, message string) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort() // Asegura que no se ejecute más lógica después.
}

// Función para agregar un Estilo a la Base de Datos
func CrearEstilo(c *gin.Context) {
	var request dto.CrearEstiloRequest

	// Validamos el Body
	if err := c.ShouldBindJSON(&request); err != nil {
		handleErrorEstilo(c, err, http.StatusBadRequest, err.Error())
		return
	}
	// Creación de Estilo
	estilo := models.Estilo{
		NombreEstilo:      request.NombreEstilo,
		DescripcionEstilo: request.DescripcionEstilo,
	}

	// Agregamos el Estilo a la base de Datos
	if err := configs.DB.Create(&estilo).Error; err != nil {
		handleErrorEstilo(c, err, http.StatusInternalServerError, "No se pudo crear el Estilo")
		return
	}

	// Retornamos mensaje con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Estilo creado con éxito"})

	// Response
	response := dto.EstiloResponse{
		ID:                estilo.ID,
		CreatedAt:         estilo.CreatedAt,
		UpdatedAt:         estilo.UpdatedAt,
		NombreEstilo:      estilo.NombreEstilo,
		DescripcionEstilo: estilo.DescripcionEstilo,
	}

	// Respondemos con Estilo Creado
	c.JSON(http.StatusCreated, gin.H{"estilo": response})
}

// Función para obtener todos los Estilos
func ObtenerEstilos(c *gin.Context) {
	var estilos []models.Estilo
	var estiloResponse []dto.EstiloResponse

	// Buscar los estilos en la base de datos
	if err := configs.DB.Where("deleted_at IS NULL").Find(&estilos).Error; err != nil {
		handleErrorEstilo(c, err, http.StatusInternalServerError, "No se pudieron obtener los Estilo")
		return
	}

	// Response
	for _, estilo := range estilos {
		estiloResponse = append(estiloResponse, dto.EstiloResponse{
			ID:                estilo.ID,
			CreatedAt:         estilo.CreatedAt,
			UpdatedAt:         estilo.UpdatedAt,
			NombreEstilo:      estilo.NombreEstilo,
			DescripcionEstilo: estilo.DescripcionEstilo,
		})
	}

	// Mostrar los Estilos
	c.JSON(http.StatusOK, gin.H{"estilos": estiloResponse})
}

// Función para obtener un estilo de acuerdo a su ID
func ObtenerEstilo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64) // Convertir ID a uint
	if err != nil {
		handleErrorEstilo(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	var estilo models.Estilo
	var estiloResponse dto.EstiloResponse

	// Buscar el Estilo por su ID
	if err := configs.DB.Where("id = ? AND deleted_at IS NULL", id).First(&estilo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			handleErrorEstilo(c, nil, http.StatusNotFound, "Estilo no encontrado")
			return
		}
		handleErrorEstilo(c, err, http.StatusInternalServerError, "Error al obtener el Estilo")
		return
	}

	// Construir la respuesta si se encuentra el estilo
	estiloResponse = dto.EstiloResponse{
		ID:                estilo.ID,
		CreatedAt:         estilo.CreatedAt,
		UpdatedAt:         estilo.UpdatedAt,
		NombreEstilo:      estilo.NombreEstilo,
		DescripcionEstilo: estilo.DescripcionEstilo,
	}

	// Enviar respuesta con los datos del estilo
	c.JSON(http.StatusOK, gin.H{"estilo": estiloResponse})
}

// Función para actualizar datos de un Estilo
func ActualizarEstilo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	var request dto.ActualizarEstiloRequest
	var estilo models.Estilo

	// Bind JSON del request a la estructura DTO
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Datos inválidos")
		return
	}

	// Buscar Estilo de acuerdo a su ID
	if err := configs.DB.First(&estilo, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Estilo no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener el Estilo")
		return
	}

	// Actualizar datos de Estilo
	estilo.NombreEstilo = request.NombreEstilo
	estilo.DescripcionEstilo = request.DescripcionEstilo

	// Guardar datos en la base de datos
	if err := configs.DB.Save(&estilo).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo actualizar el Estilo")
		return
	}

	// Response
	estiloResponse := dto.EstiloResponse{
		ID:                estilo.ID,
		NombreEstilo:      estilo.NombreEstilo,
		DescripcionEstilo: estilo.DescripcionEstilo,
	}

	// Enviar/mostrar response
	c.JSON(http.StatusOK, gin.H{"Datos de estilos": estiloResponse})
}

// Función para elimar lógicamente un Estilo
func EliminarEstilo(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, err, http.StatusBadRequest, "ID inválidos")
		return
	}

	var estilo models.Estilo

	// Buacar el Estilo por su ID
	if err := configs.DB.Unscoped().First(&estilo, id).Error; err != nil {
		HandleError(c, err, http.StatusNotFound, "Estilo no encontrado")
		return
	}

	// Verificar si el Estilo ya esta eliminado
	//if estilo.DeletedAt != nil {
	//	handleErrorEstilo(c, nil, http.StatusBadRequest, "El Estilo ya está eliminado")
	//	return
	//}
	if estilo.DeletedAt != nil && !estilo.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "El Estilo ya está eliminado")
		return
	}

	// Establecer la fecha de la eliminación lógica
	now := time.Now()
	estilo.DeletedAt = &now

	// Actualizar el registro de estilo en la base de datos
	if err := configs.DB.Save(&estilo).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo eliminar el Estilo")
		return
	}

	// Reponder/mostrar con éxito
	c.JSON(http.StatusOK, gin.H{"message": "Estilo eliminado exitosamente"})
}

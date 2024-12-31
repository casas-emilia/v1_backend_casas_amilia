package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"
	"v1_prefabricadas/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para crear Portada
func CrearPortada(c *gin.Context) {
	var request dto.CrearPortadaRequest
	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Validamos el Body
	/* if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	} */

	if err := c.ShouldBind(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Manejo de la imagen de perfil
	fileHeader, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required", "details": err.Error()})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "details": err.Error()})
		return
	}
	defer file.Close()

	// Subir la imagen a AWS S3
	url, err := services.UploadToS3(file, "portadas", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
		return
	}

	// Creación de Portada
	portada := models.Portada{
		NombrePortada: request.NombrePortada,
		//Image:         request.Image,
		Image:     url,
		EmpresaID: uint(empresaID),
	}

	// Agregamos Portada a la base de Datos
	if err := configs.DB.Create(&portada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo crear la Portada")
		return
	}

	// Response
	portadaResponse := dto.PortadaResponse{
		ID:            portada.ID,
		CreatedAt:     portada.CreatedAt,
		UpdatedAt:     portada.UpdatedAt,
		NombrePortada: portada.NombrePortada,
		Image:         portada.Image,
		EmpresaID:     portada.EmpresaID,
	}

	// Responder
	c.JSON(http.StatusCreated, gin.H{
		"message": "Portada creada con éxito",
		"portada": portadaResponse,
	})
}

// Función para obtener todas las Portadas
func ObtenerPortadas(c *gin.Context) {
	var portadas []models.Portada
	var portadasResponse []dto.PortadaResponse
	idParam := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID inválido")
		return
	}

	// Buscamos todas las portadas de la empresa
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").Find(&portadas).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Portdas no encontradas")
		return
	}

	for _, portada := range portadas {
		portadasResponse = append(portadasResponse, dto.PortadaResponse{
			ID:            portada.ID,
			CreatedAt:     portada.CreatedAt,
			UpdatedAt:     portada.UpdatedAt,
			NombrePortada: portada.NombrePortada,
			Image:         portada.Image,
			EmpresaID:     portada.EmpresaID,
		})
	}

	// Mostrar/enviar portadas
	c.JSON(http.StatusOK, gin.H{"portadas": portadasResponse})

}

// Función para obetener una portada de acuerdo al ID enviado
func ObtenerPortada(c *gin.Context) {
	/* func parseUintParam(c *gin.Context, paramName string) (uint64, error) {
		idParam := c.Param(paramName)
		return strconv.ParseUint(idParam, 10, 64)
	} */
	var portada models.Portada
	var portadaResponse dto.PortadaResponse
	idParamPortada := c.Param("portadaID")
	portadaID, err := strconv.ParseUint(idParamPortada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID portada inválido")
		return
	}

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "iD empresa inválido")
	}

	// Buscar la portada en la base de datos de acuerdo al ID enviado desde el path
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_At IS NULL").First(&portada, portadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Portada no encontrado")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Portada")
		return
	}

	portadaResponse = dto.PortadaResponse{
		ID:            portada.ID,
		CreatedAt:     portada.CreatedAt,
		UpdatedAt:     portada.UpdatedAt,
		NombrePortada: portada.NombrePortada,
		Image:         portada.Image,
	}

	// Responder/enviar Response con éxito
	c.JSON(http.StatusOK, gin.H{"portada": portadaResponse})
}

// Función para actualizar datos de una portada
func ActualizarPortada(c *gin.Context) {
	var request dto.ActualizarPortadaRequest
	var portada models.Portada
	var portadaResponse dto.PortadaResponse

	// Obtener los parámetros de la ruta
	idParamEmpresa := c.Param("empresaID")
	idParamPortada := c.Param("portadaID")

	// Convertir los parámetros a uint
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID empresa inválido")
		return
	}

	portadaID, err := strconv.ParseUint(idParamPortada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID portada inválido")
		return
	}

	// Bind JSON del request a la estructura del DTO
	if err := c.ShouldBind(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Buscar la portada en la base de datos
	if err := configs.DB.Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").First(&portada, portadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Portada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Portada")
		return
	}

	// Si se proporciona una nueva imagen, subirla a S3
	if fileHeader, err := c.FormFile("image"); err == nil {
		// Si hay una nueva imagen, subirla a S3
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image", "details": err.Error()})
			return
		}
		defer file.Close()

		// Subir la nueva imagen a S3
		url, err := services.UploadToS3(file, "portadas", fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
			return
		}
		// Actualizar la URL de la imagen en la portada
		portada.Image = url
	}

	// Actualizar otros datos de la portada
	portada.NombrePortada = request.NombrePortada

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&portada).Error; err != nil {
		HandleError(c, nil, http.StatusInternalServerError, "Error: No se pudo actualizar los datos de la Portada")
		return
	}

	// Responder con éxito y la portada actualizada
	portadaResponse = dto.PortadaResponse{
		ID:            portada.ID,
		CreatedAt:     portada.CreatedAt,
		UpdatedAt:     portada.UpdatedAt,
		NombrePortada: portada.NombrePortada,
		Image:         portada.Image,
		EmpresaID:     portada.EmpresaID,
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Portada actualizada exitosamente",
		"portada": portadaResponse,
	})
}

// Función para eliminar lógicamente una Portada de acuerdo al ID enviado
func EliminarPortada(c *gin.Context) {
	var portada models.Portada

	idParamEmpresa := c.Param("empresaID")
	idParamPortada := c.Param("portadaID")

	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	portadaID, err := strconv.ParseUint(idParamPortada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Portada inválido")
		return
	}

	// Buscar la Portada en la base de datos
	if err := configs.DB.Unscoped().Where("empresa_id = ?", empresaID).First(&portada, portadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Portada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Portada")
		return
	}

	// Verificar si la Portada ya se encuentra eliminada lógicamente
	if portada.DeletedAt != nil && !portada.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "La Portada ya se encuentra eliminada")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	portada.DeletedAt = &now

	// Guardar en la base de datos
	if err := configs.DB.Save(&portada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar la Portada")
	}

	// Mostrar mensaje de éxito de la eliminación lógica
	c.JSON(http.StatusOK, gin.H{"message": "Portada eliminada exitosamente"})
}

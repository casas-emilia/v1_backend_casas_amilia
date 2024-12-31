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

// Función para crear una imagen para la Prefabricada
func CrearImagen_prefabricada(c *gin.Context) {
	//var request dto.CrearImagen_prefabricadaRequest
	var imagen_prefabricada models.Imagen_prefabricada
	var imagen_prefabricadaResponse dto.Imagen_prefabricadaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Validamos el body
	/* if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	} */

	// Manejo de la imagen
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
	url, err := services.UploadToS3(file, "imagenes_prefabricadas", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
		return
	}

	// Crear la Imagen_prefabricada
	//imagen_prefabricada.Image = request.Image
	imagen_prefabricada.Image = url
	imagen_prefabricada.PrefabricadaID = uint(prefabricadaID)

	// Guardamos en la base de datos
	if err := configs.DB.Create(&imagen_prefabricada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No se pudo guardar la Imagen")
		return
	}

	imagen_prefabricadaResponse = dto.Imagen_prefabricadaResponse{
		ID:             imagen_prefabricada.ID,
		CreatedAt:      imagen_prefabricada.CreatedAt,
		UpdatedAt:      imagen_prefabricada.UpdatedAt,
		Image:          imagen_prefabricada.Image,
		PrefabricadaID: imagen_prefabricada.PrefabricadaID,
	}

	// Mostrar/enviar mensaje exitoso y response
	c.JSON(http.StatusOK, gin.H{"message": "Imagen guardada con éxito"})
	c.JSON(http.StatusCreated, gin.H{"imagen": imagen_prefabricadaResponse})

}

// Función para obtener todas las Imagenes de una prefabricada
func ObtenerImagenesPrefabricadas(c *gin.Context) {
	var imagenes []models.Imagen_prefabricada
	var imagenesResponse []dto.Imagen_prefabricadaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Buscar las imagenes de la prefabricada en la base de datos
	if err := configs.DB.Where("prefabricada_id = ? AND deleted_at IS NULL", prefabricadaID).Find(&imagenes).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Imagenes no encontradas")
		return
	}

	for _, imagen := range imagenes {
		imagenesResponse = append(imagenesResponse, dto.Imagen_prefabricadaResponse{
			ID:             imagen.ID,
			CreatedAt:      imagen.CreatedAt,
			UpdatedAt:      imagen.UpdatedAt,
			Image:          imagen.Image,
			PrefabricadaID: imagen.PrefabricadaID,
		})
	}

	// Mostrar/enviar imagenes_prefabricadas
	c.JSON(http.StatusOK, gin.H{"Imagenes_prefabricadas": imagenesResponse})
}

// Función para obtener una imagen de una prefabricada
func ObtenerImagePrefabricada(c *gin.Context) {
	var imagen models.Imagen_prefabricada
	var imagenesResponse dto.Imagen_prefabricadaResponse

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamImagenPrefabricada := c.Param("imagenPrefabricadaID")
	imagenPrefabricadaID, err := strconv.ParseUint(idParamImagenPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen_prefabricada inválido")
		return
	}

	// Buacar Imagen_prefabricada de acuerdo a su ID y a la Prefabricada a la que pertenece
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&imagen, imagenPrefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen_prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la imagen_prefabricada")
		return
	}

	imagenesResponse = dto.Imagen_prefabricadaResponse{
		ID:             imagen.ID,
		CreatedAt:      imagen.CreatedAt,
		UpdatedAt:      imagen.UpdatedAt,
		Image:          imagen.Image,
		PrefabricadaID: imagen.PrefabricadaID,
	}

	// Responder/enviar Imagen_prefabricada
	c.JSON(http.StatusOK, gin.H{"Imagen_prefabricada": imagenesResponse})
}

// Función para actualizar datos de una imagen_prefabricada
func ActualizarImagenPrefabricada(c *gin.Context) {
	//var request dto.ActualizarImagen_prefabricadaRequest
	var imagen models.Imagen_prefabricada
	var imagenesResponse dto.Imagen_prefabricadaResponse

	// Bind json del request a la estructura del dto
	/* if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos")
		return
	} */

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamImagenPrefabricada := c.Param("imagenPrefabricadaID")
	imagenPrefabricadaID, err := strconv.ParseUint(idParamImagenPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen_prefabricada inválido")
		return
	}

	// Buscamos la Imagen_prefabricada en la base de datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").First(&imagen, imagenPrefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen_prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la imagen_prefabricada")
		return
	}

	// actualizar los datos de la Imagen_prefabricada
	//imagen.Image = request.Image

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
		url, err := services.UploadToS3(file, "imagenes_prefabricadas", fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
			return
		}
		// Actualizar la URL de la imagen
		imagen.Image = url
	}

	// Guardar cambios en la base de datos
	if err := configs.DB.Save(&imagen).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo guardar los cambios")
		return
	}

	imagenesResponse = dto.Imagen_prefabricadaResponse{
		ID:             imagen.ID,
		CreatedAt:      imagen.CreatedAt,
		UpdatedAt:      imagen.UpdatedAt,
		Image:          imagen.Image,
		PrefabricadaID: imagen.PrefabricadaID,
	}

	// Responder/enviar un mensaje de éxito y el response de la Imagen_prefabricada
	c.JSON(http.StatusOK, gin.H{
		"message":             "Cambios guardados con éxito",
		"Imagen_prefabricada": imagenesResponse,
	})

}

// Función para eliminar lógicamente una imagen_prefabricada
func EliminarImagenPrefabricada(c *gin.Context) {
	var imagen models.Imagen_prefabricada

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	idParamImagenPrefabricada := c.Param("imagenPrefabricadaID")
	imagenPrefabricadaID, err := strconv.ParseUint(idParamImagenPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen_prefabricada inválido")
		return
	}

	// Buscar la imagen_prefabricada en la base datos
	if err := configs.DB.Where("prefabricada_id = ?", prefabricadaID).Where("deleted_at IS NULL").Unscoped().First(&imagen, imagenPrefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen_prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la imagen_prefabricada")
		return
	}

	// Verificar si la Imagen_prefabricada ya se encuentra eliminada lógicamente
	if imagen.DeletedAt != nil && !imagen.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "La imagen ya se encuentra eliminada")
		return
	}

	// Poner fecha y hora de la eliminación
	now := time.Now()
	imagen.DeletedAt = &now

	// Guardar fecha y hora de la eliminación lógica en la base de datos
	if err := configs.DB.Save(&imagen).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar la Imagen")
		return
	}

	// Mostrar/enviar un mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Imagen eliminada exitosamente"})
}

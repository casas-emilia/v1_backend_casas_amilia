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

// Función para Crear Una Imagen de Noticia
func CrearImagenNoticia(c *gin.Context) {
	//var request dto.CrearImagen_noticiaRequest
	var imagen models.Imagen_noticia
	var imagenResponse dto.Imagen_noticiaResponse

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	// Validamos el body
	// if err := c.ShouldBindJSON(&request); err != nil {
	// 	HandleError(c, err, http.StatusInternalServerError, "Error de datos")
	// 	return
	// }

	// // Creamos la Imagen
	// imagen.Image = request.Image
	//imagen.NoticiaID = uint(noticiaID)

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
	url, err := services.UploadToS3(file, "noticias", fileHeader.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
		return
	}

	imagen.NoticiaID = uint(noticiaID)
	imagen.Image = url

	// Guardar la Imagen en la base da datos
	if err := configs.DB.Create(&imagen).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo guardar la Imagen")
		return
	}

	imagenResponse = dto.Imagen_noticiaResponse{
		ID:        imagen.ID,
		CreatedAt: imagen.CreatedAt,
		UpdatedAt: imagen.UpdatedAt,
		Image:     imagen.Image,
		NoticiaID: imagen.NoticiaID,
	}

	// Mostrar/enviar mensaje exitoso e Imagen
	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen guardada con éxito",
	})
	c.JSON(http.StatusCreated, gin.H{
		"imagen": imagenResponse,
	})

}

// Obtener Todas las Imagenes de una Noticia
func ObtenerImagenesNoticias(c *gin.Context) {
	var imagenes []models.Imagen_noticia
	var imagenesResponse []dto.Imagen_noticiaResponse

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	// Buscar todas las imagenes de una noticia
	if err := configs.DB.Where("noticia_id = ?", noticiaID).Where("deleted_at IS NULL").Find(&imagenes).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener las imagenes de la Noticia")
		return
	}

	// Si no hay imágenes, devolver un array vacío con status 200
	if len(imagenes) == 0 {
		// Devolver array vacío en lugar de error 404
		c.JSON(http.StatusOK, gin.H{
			"imagenes_noticia": []interface{}{},
		})
		return
	}

	for _, imagen := range imagenes {
		imagenesResponse = append(imagenesResponse, dto.Imagen_noticiaResponse{
			ID:        imagen.ID,
			CreatedAt: imagen.CreatedAt,
			UpdatedAt: imagen.UpdatedAt,
			Image:     imagen.Image,
			NoticiaID: imagen.NoticiaID,
		})
	}

	// Mostrat/enviar imagenes de la Noticia
	c.JSON(http.StatusOK, gin.H{
		"imagenes_noticia": imagenesResponse,
	})

}

// Obtener una imagen de acuerdo a su ID
func ObtenerImagenNoticia(c *gin.Context) {
	var imagen models.Imagen_noticia
	var imagenResponse dto.Imagen_noticiaResponse

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	idParamImagen := c.Param("imagenNoticiaID")
	imagenNoticiaID, err := strconv.ParseUint(idParamImagen, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen inválido")
		return
	}

	// Buscar Imagen en la base de datos de acuerdo a su ID
	if err := configs.DB.Where("noticia_id = ? AND id = ?", noticiaID, imagenNoticiaID).Where("deleted_at IS NULL").First(&imagen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Imagen")
		return
	}

	imagenResponse = dto.Imagen_noticiaResponse{
		ID:        imagen.ID,
		CreatedAt: imagen.CreatedAt,
		UpdatedAt: imagen.UpdatedAt,
		Image:     imagen.Image,
		NoticiaID: imagen.NoticiaID,
	}

	// Mostrar/enviar imagen
	c.JSON(http.StatusOK, gin.H{
		"imagen_noticia": imagenResponse,
	})
}

// Función para actualizar una imagen
func ActualizarImagenNoticia(c *gin.Context) {
	var imagen models.Imagen_noticia
	var imagenResponse dto.Imagen_noticiaResponse

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	idParamImagen := c.Param("imagenNoticiaID")
	imagenNoticiaID, err := strconv.ParseUint(idParamImagen, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen inválido")
		return
	}

	// Buscar la Imagen en la Base de datos
	if err := configs.DB.Where("noticia_id = ? AND id = ?", noticiaID, imagenNoticiaID).Where("deleted_at IS NULL").First(&imagen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener Imagen")
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
		url, err := services.UploadToS3(file, "noticias", fileHeader.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to S3", "details": err.Error()})
			return
		}
		// Actualizar la URL de la imagen
		imagen.Image = url
	}

	// guardar cambios en la base de datos
	if err := configs.DB.Save(&imagen).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo actualizar la imagen")
		return
	}

	imagenResponse = dto.Imagen_noticiaResponse{
		ID:        imagen.ID,
		CreatedAt: imagen.CreatedAt,
		UpdatedAt: imagen.UpdatedAt,
		Image:     imagen.Image,
		NoticiaID: imagen.NoticiaID,
	}

	// Mostrar/enviar mensaje de actualización exitosa y la Imagen
	c.JSON(http.StatusOK, gin.H{
		"message":        "Imagen actualizada con éxito",
		"imagen_noticia": imagenResponse,
	})
}

// Función para eliminar lógicamente una imagen
func EliminarImagenNoticia(c *gin.Context) {
	var imagen models.Imagen_noticia

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	idParamImagen := c.Param("imagenNoticiaID")
	imagenNoticiaID, err := strconv.ParseUint(idParamImagen, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Imagen inválido")
		return
	}

	// Buacar Imagen en la Base de datos
	if err := configs.DB.Where("noticia_id = ? AND id = ?", noticiaID, imagenNoticiaID).Where("deleted_at IS NULL").First(&imagen).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Imagen no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de buscar la Imagen")
		return
	}

	if imagen.DeletedAt != nil && imagen.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "Imagen ya se encuentra eliminada")
		return
	}

	// Poner fecha y hora de la eliminación lógica
	now := time.Now()
	imagen.DeletedAt = &now

	// guardar Fecha y hora de la eliminación lógica
	if err := configs.DB.Save(&imagen).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de eliminar la Imagen")
		return
	}

	// Mostrar/enviar un mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Imagen eliminada con éxito",
	})

}

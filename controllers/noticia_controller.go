package controllers

import (
	"errors"
	"net/http"
	"strconv"
	"time"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/dto"
	"v1_prefabricadas/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Función para crear un Noticia
func CrearNoticia(c *gin.Context) {
	var request dto.CrearNoticiaRequest
	var noticia models.Noticia
	var noticiaResponse dto.NoticiaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido ")
		return
	}

	// usuarioID, err := helpers.ValidarUsuarioID(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	// Bind JSON de request a la estructura de dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Creamos la Noticia
	noticia.TituloNoticia = request.TituloNoticia
	noticia.DesarrolloNoticia = request.DesarrolloNoticia
	//noticia.UsuarioID = uint(usuarioID)
	//noticia.UsuarioID = usuarioID
	noticia.EmpresaID = uint(empresaID)

	// Guardamos la noticia en la base de datos
	if err := configs.DB.Create(&noticia).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo crear la Noticia")
		return
	}

	noticiaResponse = dto.NoticiaResponse{
		ID:                noticia.ID,
		CreatedAt:         noticia.CreatedAt,
		UpdatedAt:         noticia.UpdatedAt,
		TituloNoticia:     noticia.TituloNoticia,
		DesarrolloNoticia: noticia.DesarrolloNoticia,
		EmpresaID:         noticia.EmpresaID,
	}

	// Mostrar/enviar mensaje exitoso y noticiaResponse
	c.JSON(http.StatusOK, gin.H{
		"message": "Noticia creada con éxito",
		"noticia": noticiaResponse,
	})

}

// Función para obtener todas las noticias de una empresa
func ObtenerNoticiasEmpresa(c *gin.Context) {
	var noticias []models.Noticia
	var noticiasResponse []dto.NoticiaResponse

	// Parámetros de paginación desde la solicitud
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 12
	}

	offset := (page - 1) * limit

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Calcular el total de noticias sin paginación
	var totalCount int64
	if err := configs.DB.
		Model(&models.Noticia{}). // Cambié `Noticia{}` a `models.Noticia{}`
		Where("deleted_at IS NULL").
		Where("empresa_id = ?", empresaID).
		Count(&totalCount).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al contar las Noticias/Actividas")
		return
	}

	// Buscar todas las noticias con paginación
	if err := configs.DB.
		Where("deleted_at IS NULL").
		Where("empresa_id = ?", empresaID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&noticias).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener las Noticias/Actividas")
		return
	}

	// Verificar si no hay noticias
	if len(noticias) == 0 {
		HandleError(c, nil, http.StatusNotFound, "Sin Noticias/Actividades por el momento")
		return
	}

	// Crear la respuesta
	for _, noticia := range noticias {
		noticiasResponse = append(noticiasResponse, dto.NoticiaResponse{
			ID:                noticia.ID,
			CreatedAt:         noticia.CreatedAt,
			UpdatedAt:         noticia.UpdatedAt,
			TituloNoticia:     noticia.TituloNoticia,
			DesarrolloNoticia: noticia.DesarrolloNoticia,
			EmpresaID:         noticia.EmpresaID,
		})
	}

	// Retornar las noticias junto con información de paginación
	c.JSON(http.StatusOK, gin.H{
		"noticias": noticiasResponse, // Enviar el array `noticiasResponse` con el DTO
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": totalCount,
		},
	})
}

// Función para obterner una noticia especifica de una Empresa
func ObtenerNoticiaEmpresa(c *gin.Context) {
	var noticia models.Noticia
	var noticiaResponse dto.NoticiaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia no válido")
		return
	}

	// Buscar en la base de datos noticia de acuerdo al id de la empresa y al id de la noticia
	if err := configs.DB.Where("deleted_at IS NULL").Where("empresa_id = ? AND id = ?", empresaID, noticiaID).First(&noticia).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Noticia no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener la Noticia")
		return
	}

	noticiaResponse = dto.NoticiaResponse{
		ID:                noticia.ID,
		CreatedAt:         noticia.CreatedAt,
		UpdatedAt:         noticia.UpdatedAt,
		TituloNoticia:     noticia.TituloNoticia,
		DesarrolloNoticia: noticia.DesarrolloNoticia,
		EmpresaID:         noticia.EmpresaID,
	}

	// Mostar/enviar Noticia
	c.JSON(http.StatusOK, gin.H{
		"noticia": noticiaResponse,
	})
}

// Función para obtener todas las Noticias de cada Usuario(por su ID)
// func ObtenerNoticiasUsuarios(c *gin.Context) {
// 	var noticias []models.Noticia
// 	var noticiasResponse []dto.NoticiaResponse

// 	idParamUsuario := c.Param("usuarioID")
// 	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
// 	if err != nil {
// 		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
// 		return
// 	}

// 	// Buscar todas las Noticias de cada usuario
// 	if err := configs.DB.Where("usuario_id = ?", usuarioID).Where("deleted_at IS NULL").Find(&noticias).Error; err != nil {
// 		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de obtener las Noticias/Actividades")
// 		return
// 	}

// 	if len(noticias) == 0 {
// 		HandleError(c, nil, http.StatusNotFound, "Usuario sin noticias/actividades creadas aún")
// 		return
// 	}

// 	for _, noticia := range noticias {
// 		noticiasResponse = append(noticiasResponse, dto.NoticiaResponse{
// 			ID:                noticia.ID,
// 			CreatedAt:         noticia.CreatedAt,
// 			UpdatedAt:         noticia.UpdatedAt,
// 			TituloNoticia:     noticia.TituloNoticia,
// 			DesarrolloNoticia: noticia.DesarrolloNoticia,
// 			EmpresaID:         noticia.EmpresaID,
// 		})
// 	}

// 	// Mostrar/enviar Noticias creadas por el Usuario
// 	c.JSON(http.StatusOK, gin.H{
// 		"noticias": noticiasResponse,
// 	})
// }

// Función para obtener una Noticia de Usuario
// func ObtenerNoticiaUsuario(c *gin.Context) {
// 	var noticia models.Noticia
// 	var noticiaResponse dto.NoticiaResponse

// 	idParamUsuario := c.Param("usuarioID")
// 	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
// 	if err != nil {
// 		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
// 		return
// 	}

// 	idParamNoticia := c.Param("noticiaID")
// 	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
// 	if err != nil {
// 		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
// 		return
// 	}

// 	// Buscar Noticia en la Base de datos de acuerdo al ID de noticia y al ID del Usuario
// 	if err := configs.DB.Where("usuario_id = ? AND id = ?", usuarioID, noticiaID).Where("deleted_at IS NULL").First(&noticia).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			HandleError(c, nil, http.StatusNotFound, "Noticia no encontrada")
// 			return
// 		}
// 		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Noticia")
// 		return
// 	}

// 	noticiaResponse = dto.NoticiaResponse{
// 		ID:                noticia.ID,
// 		CreatedAt:         noticia.CreatedAt,
// 		UpdatedAt:         noticia.UpdatedAt,
// 		TituloNoticia:     noticia.TituloNoticia,
// 		DesarrolloNoticia: noticia.DesarrolloNoticia,
// 		EmpresaID:         noticia.EmpresaID,
// 	}

// 	// Mostrar/enviar Noticia
// 	c.JSON(http.StatusOK, gin.H{
// 		"noticia": noticiaResponse,
// 	})
// }

// Función para actualizar Noticia
func ActualizarNoticia(c *gin.Context) {
	var request dto.ActualizarNoticiaRequest
	var noticia models.Noticia
	var noticiaResponse dto.NoticiaResponse

	// idParamUsuario := c.Param("usuarioID")
	// usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	// if err != nil {
	// 	HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
	// 	return
	// }

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	// validamos request
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, "Error de datos "+err.Error())
		return
	}

	// Buacamos la noticia en la bade de datos
	if err := configs.DB.Where("id = ?", noticiaID).Where("deleted_at IS NULL").First(&noticia).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Noticia no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de buscar los datos de la Noticia")
		return
	}

	// Actualizar los datos de la Noticia
	if request.TituloNoticia != "" {
		noticia.TituloNoticia = request.TituloNoticia
	} else {
		HandleError(c, nil, http.StatusBadRequest, "El Título de la Noticia no debe estar vacio")
		return
	}

	if request.DesarrolloNoticia != "" {
		noticia.DesarrolloNoticia = request.DesarrolloNoticia
	} else {
		HandleError(c, nil, http.StatusBadRequest, "El Desarrollo de la Noticia no debe estar vacio")
		return
	}

	// Guardar en la base de datos
	if err := configs.DB.Save(&noticia).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de actualizar los datos de la Noticia")
		return
	}

	noticiaResponse = dto.NoticiaResponse{
		ID:                noticia.ID,
		CreatedAt:         noticia.CreatedAt,
		UpdatedAt:         noticia.UpdatedAt,
		TituloNoticia:     noticia.TituloNoticia,
		DesarrolloNoticia: noticia.DesarrolloNoticia,
		EmpresaID:         noticia.EmpresaID,
	}

	// Mostrar/enviar NoticiaResponse y un mensaje de actualización exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Datos actualizados exitosamente",
		"noticia": noticiaResponse,
	})
}

// Eliminar Noticia lógicamente
func EliminarNoticia(c *gin.Context) {
	var noticia models.Noticia

	/* idParamUsuario := c.Param("usuarioID")
	usuarioID, err := strconv.ParseUint(idParamUsuario, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Usuario inválido")
		return
	} */

	// usuarioID, err := helpers.ValidarUsuarioID(c)
	// if err != nil {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	// 	return
	// }

	idParamNoticia := c.Param("noticiaID")
	noticiaID, err := strconv.ParseUint(idParamNoticia, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Noticia inválido")
		return
	}

	// Buscar la Noticia en la base de datos de acuerdo al ID de la Noticia y al ID del Usuario que la Creo
	if err := configs.DB.Where("id = ?", noticiaID).Where("deleted_at IS NULL").First(&noticia).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Noticia no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al tratar de buscar la Noticia")
		return
	}

	// Poner Fecha y Hora de la eliminación lógica
	now := time.Now()
	noticia.DeletedAt = &now

	// Guardar en la base da datos la fecha y hora de la eliminación lógica
	if err := configs.DB.Save(&noticia).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar la Noticia")
		return
	}

	// Enviar/mostrar mensaje de eliminación exitosa
	c.JSON(http.StatusOK, gin.H{
		"message": "Noticia eliminada exitosamente",
	})

}

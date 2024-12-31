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

func CrearPrefabricada(c *gin.Context) {
	var prefabricada models.Prefabricada
	var request dto.CrearPrefabricadaRequest
	var prefabricadaResponse dto.PrefabricadaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID de Empresa inválido")
		return
	}

	// Validamos el body
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, err, http.StatusBadRequest, err.Error())
		return
	}

	// Crear la Prefabricada
	prefabricada.NombrePrefabricada = request.NombrePrefabricada
	prefabricada.M2 = request.M2
	prefabricada.Garantia = request.Garantia
	prefabricada.Eslogan = request.Eslogan
	prefabricada.Descripcion = request.Descripcion
	prefabricada.Destacada = request.Destacada
	prefabricada.Oferta = request.Oferta
	prefabricada.CategoriaID = request.CategoriaID
	prefabricada.EmpresaID = uint(empresaID)
	prefabricada.EstiloID = request.EstiloID
	prefabricada.TipoID = request.TipoID

	// Guardar en la base de datos la Prefabricada
	if err := configs.DB.Create(&prefabricada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo crear la Prefabricada")
	}

	// Response
	prefabricadaResponse = dto.PrefabricadaResponse{
		ID:                 prefabricada.ID,
		NombrePrefabricada: prefabricada.NombrePrefabricada,
		M2:                 prefabricada.M2,
		Garantia:           prefabricada.Garantia,
		Eslogan:            prefabricada.Eslogan,
		Descripcion:        prefabricada.Descripcion,
		Destacada:          prefabricada.Destacada,
		Oferta:             prefabricada.Oferta,
		CategoriaID:        prefabricada.CategoriaID,
		EmpresaID:          prefabricada.EmpresaID,
		EstiloID:           prefabricada.EstiloID,
		TipoID:             prefabricada.TipoID,
	}

	c.JSON(http.StatusOK, gin.H{"prefabricada": prefabricadaResponse})
}

// Función para obtener todas las Prefabricadas con paginación
func ObtenerPrefabricadas(c *gin.Context) {
	var prefabricadas []models.Prefabricada
	var prefabricadasResponse []dto.PrefabricadaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	// Obtener filtros opcionales
	categoriaID := c.Query("categoria_id")
	tipoID := c.Query("tipo_id")
	destacada := c.Query("destacada")
	oferta := c.Query("oferta")

	// Obtener parámetros de paginación
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1 // Si no se pasa una página válida, por defecto será la página 1
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "12")) // Default: 12 resultados por página
	if err != nil || limit < 1 {
		limit = 12 // Si no se pasa un límite válido, por defecto serán 12 resultados
	}

	// Calcular el offset (paginación)
	offset := (page - 1) * limit

	// Iniciar consulta base
	query := configs.DB.
		Preload("Imagen_prefabricada", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar imágenes eliminadas lógicamente
		}).
		Preload("Caracteristica", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar características eliminadas lógicamente
		}).
		Preload("Precio", func(db *gorm.DB) *gorm.DB {
			return db.Where("precios.deleted_at IS NULL") // Filtra precios no eliminados
		}).Preload("Precio.Incluye", func(db *gorm.DB) *gorm.DB {
		return db.Where("incluyes.deleted_at IS NULL") // Filtra los "incluye" no eliminados
	}).
		Where("prefabricadas.empresa_id = ?", empresaID).
		Where("prefabricadas.deleted_at IS NULL")

	// Aplicar filtro por categoria_id si está presente
	if categoriaID != "" {
		query = query.Where("prefabricadas.categoria_id = ?", categoriaID)
	}

	// Aplicar filtro por tipo_id si está presente
	if tipoID != "" {
		query = query.Joins("JOIN tipos ON tipos.id = prefabricadas.tipo_id").
			Where("tipos.id = ?", tipoID)
	}

	// Aplicar filtro por destacada
	if destacada != "" {
		query = query.Where("prefabricadas.destacada = ?", destacada)
	}

	// Aplicar filtro por Oferta
	if oferta != "" {
		query = query.Where("prefabricadas.oferta = ?", oferta)
	}

	// Ejecutar consulta con paginación
	if err := query.Limit(limit).Offset(offset).Find(&prefabricadas).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Prefabricadas no encontradas")
		return
	}

	// Verificar si no hay prefabricadas y devolver un array vacío
	if len(prefabricadas) == 0 {
		prefabricadasResponse = []dto.PrefabricadaResponse{}
	} else {
		// Crear response de prefabricadas
		for _, prefabricada := range prefabricadas {
			var imagenes_prefabricadasResponse []dto.Imagen_prefabricadaResponse
			var caracteristicasResponse []dto.CaracteristicaResponse
			var preciosResponse []dto.PrecioResponse

			// Convertir imágenes de prefabricada a DTO
			for _, imagen := range prefabricada.Imagen_prefabricada {
				imagenes_prefabricadasResponse = append(imagenes_prefabricadasResponse, dto.Imagen_prefabricadaResponse{
					ID:             imagen.ID,
					CreatedAt:      imagen.CreatedAt,
					UpdatedAt:      imagen.UpdatedAt,
					Image:          imagen.Image,
					PrefabricadaID: imagen.PrefabricadaID,
				})
			}

			// Convertir características a DTO
			for _, caracteristica := range prefabricada.Caracteristica {
				caracteristicasResponse = append(caracteristicasResponse, dto.CaracteristicaResponse{
					ID:             caracteristica.ID,
					CreatedAt:      caracteristica.CreatedAt,
					UpdatedAt:      caracteristica.UpdatedAt,
					Clave:          caracteristica.Clave,
					Valor:          caracteristica.Valor,
					PrefabricadaID: caracteristica.PrefabricadaID,
				})
			}

			// Convertir precios e incluyes a DTO
			for _, precio := range prefabricada.Precio {
				var incluyesResponse []dto.IncluyeResponse

				for _, incluye := range precio.Incluye {
					incluyesResponse = append(incluyesResponse, dto.IncluyeResponse{
						ID:            incluye.ID,
						CreatedAt:     incluye.CreatedAt,
						UpdatedAt:     incluye.UpdatedAt,
						NombreIncluye: incluye.NombreIncluye,
						PrecioID:      incluye.PrecioID,
					})
				}

				preciosResponse = append(preciosResponse, dto.PrecioResponse{
					ID:                precio.ID,
					CreatedAt:         precio.CreatedAt,
					UpdatedAt:         precio.UpdatedAt,
					NombrePrecio:      precio.NombrePrecio,
					DescripcionPrecio: precio.DescripcionPrecio,
					ValorPrefabricada: precio.ValorPrefabricada,
					PrefabricadaID:    precio.PrefabricadaID,
					Incluyes:          incluyesResponse,
				})
			}

			// Añadir prefabricada a la respuesta
			prefabricadasResponse = append(prefabricadasResponse, dto.PrefabricadaResponse{
				ID:                    prefabricada.ID,
				CreatedAt:             prefabricada.CreatedAt,
				UpdatedAt:             prefabricada.UpdatedAt,
				NombrePrefabricada:    prefabricada.NombrePrefabricada,
				M2:                    prefabricada.M2,
				Garantia:              prefabricada.Garantia,
				Eslogan:               prefabricada.Eslogan,
				Descripcion:           prefabricada.Descripcion,
				Destacada:             prefabricada.Destacada,
				Oferta:                prefabricada.Oferta,
				CategoriaID:           prefabricada.CategoriaID,
				EmpresaID:             prefabricada.EmpresaID,
				EstiloID:              prefabricada.EstiloID,
				TipoID:                prefabricada.TipoID,
				ImagenesPrefabricadas: imagenes_prefabricadasResponse,
				Caracteristicas:       caracteristicasResponse,
				Precios:               preciosResponse,
			})
		}
	}

	// Mostrar/enviar response de Prefabricada con información de paginación
	c.JSON(http.StatusOK, gin.H{
		"prefabricadas": prefabricadasResponse, // Aquí ya será un array vacío en caso de no haber resultados
		"page":          page,
		"limit":         limit,
	})
}

// Función para obtener una Prefabricada de acuerdo al ID enviado
func ObtenerPrefabricada(c *gin.Context) {
	var prefabricada models.Prefabricada
	var prefabricadaResponse dto.PrefabricadaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	if err := configs.DB.
		Preload("Imagen_prefabricada", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar imágenes eliminadas lógicamente
		}).
		Preload("Caracteristica", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL") // Condición para no cargar características eliminadas lógicamente
		}).
		Preload("Precio", func(db *gorm.DB) *gorm.DB {
			return db.Where("precios.deleted_at IS NULL") // Filtra precios no eliminados
		}).Preload("Precio.Incluye", func(db *gorm.DB) *gorm.DB {
		return db.Where("incluyes.deleted_at IS NULL") // Filtra los "incluye" no eliminados
	}).
		Where("empresa_id = ?", empresaID).
		Where("deleted_at IS NULL").
		First(&prefabricada, prefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Prefabricada")
		return
	}

	var imagenes_prefabricadasResponse []dto.Imagen_prefabricadaResponse
	var caracteristicasResponse []dto.CaracteristicaResponse
	var preciosResponse []dto.PrecioResponse

	for _, imagen := range prefabricada.Imagen_prefabricada {
		imagenes_prefabricadasResponse = append(imagenes_prefabricadasResponse, dto.Imagen_prefabricadaResponse{
			ID:             imagen.ID,
			CreatedAt:      imagen.CreatedAt,
			UpdatedAt:      imagen.UpdatedAt,
			Image:          imagen.Image,
			PrefabricadaID: imagen.PrefabricadaID,
		})
	}

	for _, caracteristica := range prefabricada.Caracteristica {
		caracteristicasResponse = append(caracteristicasResponse, dto.CaracteristicaResponse{
			ID:             caracteristica.ID,
			CreatedAt:      caracteristica.CreatedAt,
			UpdatedAt:      caracteristica.UpdatedAt,
			Clave:          caracteristica.Clave,
			Valor:          caracteristica.Valor,
			PrefabricadaID: caracteristica.PrefabricadaID,
		})
	}

	for _, precio := range prefabricada.Precio {
		var incluyesResponse []dto.IncluyeResponse

		for _, incluye := range precio.Incluye {
			incluyesResponse = append(incluyesResponse, dto.IncluyeResponse{
				ID:            incluye.ID,
				CreatedAt:     incluye.CreatedAt,
				UpdatedAt:     incluye.UpdatedAt,
				NombreIncluye: incluye.NombreIncluye,
				PrecioID:      incluye.PrecioID,
			})
		}

		preciosResponse = append(preciosResponse, dto.PrecioResponse{
			ID:                precio.ID,
			CreatedAt:         precio.CreatedAt,
			UpdatedAt:         precio.UpdatedAt,
			NombrePrecio:      precio.NombrePrecio,
			DescripcionPrecio: precio.DescripcionPrecio,
			ValorPrefabricada: precio.ValorPrefabricada,
			PrefabricadaID:    precio.PrefabricadaID,
			Incluyes:          incluyesResponse,
		})
	}

	prefabricadaResponse = dto.PrefabricadaResponse{
		ID:                    prefabricada.ID,
		CreatedAt:             prefabricada.CreatedAt,
		UpdatedAt:             prefabricada.UpdatedAt,
		NombrePrefabricada:    prefabricada.NombrePrefabricada,
		M2:                    prefabricada.M2,
		Garantia:              prefabricada.Garantia,
		Eslogan:               prefabricada.Eslogan,
		Descripcion:           prefabricada.Descripcion,
		Destacada:             prefabricada.Destacada,
		Oferta:                prefabricada.Oferta,
		CategoriaID:           prefabricada.CategoriaID,
		EmpresaID:             prefabricada.EmpresaID,
		EstiloID:              prefabricada.EstiloID,
		TipoID:                prefabricada.TipoID,
		ImagenesPrefabricadas: imagenes_prefabricadasResponse,
		Caracteristicas:       caracteristicasResponse,
		Precios:               preciosResponse,
	}

	// Mostrar/enviar response
	c.JSON(http.StatusOK, gin.H{"prefabricada": prefabricadaResponse})
}

// Función para actualizar datos basicos de una Prefabricada
func ActualizarPrefabricada(c *gin.Context) {
	var request dto.ActualizarPrefabricadaRequest
	var prefabricada models.Prefabricada
	var prefabricadaResponse dto.PrefabricadaResponse

	idParamEmpresa := c.Param("empresaID")
	empresaId, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Bind json del request a la estructura del dto
	if err := c.ShouldBindJSON(&request); err != nil {
		HandleError(c, nil, http.StatusBadRequest, "Error de datos")
		return
	}

	// Buscar datos de Prefabricadas de acuerdo al ID en la base da datos
	if err := configs.DB.Where("empresa_id = ?", empresaId).Where("deleted_At IS NULL").First(&prefabricada, prefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Prefabricada")
		return
	}

	// Actualizar datos de la Prefabricada
	prefabricada.NombrePrefabricada = request.NombrePrefabricada
	prefabricada.M2 = request.M2
	prefabricada.Garantia = request.Garantia
	prefabricada.Eslogan = request.Eslogan
	prefabricada.Descripcion = request.Descripcion
	prefabricada.Destacada = request.Destacada
	prefabricada.Oferta = request.Oferta
	prefabricada.CategoriaID = request.CategoriaID
	prefabricada.EstiloID = request.EstiloID
	prefabricada.TipoID = request.TipoID

	// Guardar los cambios en la base de datos
	if err := configs.DB.Save(&prefabricada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "No pudo actualizar datos de Prefabricada")
		return
	}

	prefabricadaResponse = dto.PrefabricadaResponse{
		ID:                 prefabricada.ID,
		CreatedAt:          prefabricada.CreatedAt,
		UpdatedAt:          prefabricada.UpdatedAt,
		NombrePrefabricada: prefabricada.NombrePrefabricada,
		M2:                 prefabricada.M2,
		Garantia:           prefabricada.Garantia,
		Eslogan:            prefabricada.Eslogan,
		Descripcion:        prefabricada.Descripcion,
		Destacada:          prefabricada.Destacada,
		Oferta:             prefabricada.Oferta,
		CategoriaID:        prefabricada.CategoriaID,
		EstiloID:           prefabricada.EstiloID,
		TipoID:             prefabricada.TipoID,
	}

	// Mostrar/enviar mensaje de éxto
	c.JSON(http.StatusOK, gin.H{
		"message":      "Datos actualizados exitosamente",
		"prefabricada": prefabricadaResponse,
	})
}

// Función para eliminar una Prefabricada
func EliminarPrefabricada(c *gin.Context) {
	var prefabricada models.Prefabricada

	idParamEmpresa := c.Param("empresaID")
	empresaID, err := strconv.ParseUint(idParamEmpresa, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Empresa inválido")
		return
	}

	idParamPrefabricada := c.Param("prefabricadaID")
	prefabricadaID, err := strconv.ParseUint(idParamPrefabricada, 10, 64)
	if err != nil {
		HandleError(c, nil, http.StatusBadRequest, "ID Prefabricada inválido")
		return
	}

	// Buscar la Prefabricada en la Base de datos con su ID
	if err := configs.DB.Unscoped().Where("empresa_id = ?", empresaID).Where("deleted_at IS NULL").First(&prefabricada, prefabricadaID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			HandleError(c, nil, http.StatusNotFound, "Prefabricada no encontrada")
			return
		}
		HandleError(c, err, http.StatusInternalServerError, "Error al obtener datos de la Prefabricada")
		return
	}

	// Verificar si la Prefabricada ya se encuentra eliminada lógicamente
	if prefabricada.DeletedAt != nil && !prefabricada.DeletedAt.IsZero() {
		HandleError(c, nil, http.StatusBadRequest, "La prefabricada ya se encuentra eliminada")
		return
	}

	// Poner facha y hora de la eliminación lógica
	now := time.Now()
	prefabricada.DeletedAt = &now

	// Guardar facha y hora de la eliminación lógica de la Prefabricada en la base de datos
	if err := configs.DB.Save(&prefabricada).Error; err != nil {
		HandleError(c, err, http.StatusInternalServerError, "Error, no se pudo eliminar la Prefabricada")
		return
	}

	// Mostrar/enviar mensaje de eliminación lógica exitosa
	c.JSON(http.StatusOK, gin.H{"message": "Prefabricada eliminada exitosamente"})
}

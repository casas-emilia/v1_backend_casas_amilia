package main

import (
	"log"
	"v1_prefabricadas/configs"
	"v1_prefabricadas/models"
)

func init() {
	configs.ConnectToDB()
}

func main() {
	log.Println("Iniciando migraciones...")

	err := configs.DB.AutoMigrate(
		&models.Caracteristica{},
		&models.Categoria{},
		&models.Contacto{},
		&models.Credencial{},
		&models.Empresa{},
		&models.Estilo{},
		&models.Imagen_noticia{},
		&models.Imagen_prefabricada{},
		&models.Incluye{},
		&models.Noticia{},
		&models.Portada{},
		&models.Precio{},
		&models.Prefabricada{},
		&models.Red{},
		&models.Rol{},
		&models.Rol_usuario{},
		&models.Servicio{},
		&models.Tipo{},
		&models.Tipo_categoria{},
		&models.Usuario{},
		&models.Recuperacion{},
	)
	if err != nil {
		log.Fatalf("Error durante la migración: %v", err)
	}

	log.Println("Migraciones completadas exitosamente")

	// Seed initial data
	seedInitialData()
}

func seedInitialData() {
	// Verificar si los roles necesarios ya están creados
	for _, id := range []uint{1, 2, 3} {
		var rol models.Rol
		if err := configs.DB.First(&rol, id).Error; err != nil {
			// Si no existe el rol, crearlo según el ID
			nombreRol := ""
			descripcionRol := ""
			switch id {
			case 1:
				nombreRol = "super_administrador"
				descripcionRol = "Super administrador del sistema"
			case 2:
				nombreRol = "administrador"
				descripcionRol = "Administrador del sistema con permisos avanzados"
			case 3:
				nombreRol = "ejecutivo_ventas"
				descripcionRol = "Responsable de gestionar las ventas"
			}

			rol = models.Rol{
				ID:             id,
				NombreRol:      nombreRol,
				DescripcionRol: descripcionRol,
			}

			if result := configs.DB.Create(&rol); result.Error != nil {
				log.Fatalf("Error al crear el rol con ID %d: %v", id, result.Error)
			}
		}
	}

	log.Println("Roles iniciales verificados y creados correctamente")

	// Verificar si la empresa con ID 1 ya está creada
	var empresa models.Empresa
	if err := configs.DB.First(&empresa, 1).Error; err != nil {
		// Si no existe, crear la empresa con ID 1
		empresa = models.Empresa{
			ID:                 1,
			NombreEmpresa:      "Empresa Demo",
			DescripcionEmpresa: "Descripción de la empresa demo",
			HistoriaEmpresa:    "Historia de la empresa demo",
			MisionEmpresa:      "Misión de la empresa demo",
			VisionEmpresa:      "Visión de la empresa demo",
			UbicacionEmpresa:   "Ubicación demo",
			CelularEmpresa:     "123456789",
			EmailEmpresa:       "empresa@demo.com",
		}
		if result := configs.DB.Create(&empresa); result.Error != nil {
			log.Fatalf("Error al crear la empresa inicial: %v", result.Error)
		}
	}

	// Verificar si el usuario con ID 1 ya está creado
	var superAdminUser models.Usuario
	if err := configs.DB.First(&superAdminUser, 1).Error; err != nil {
		// Si no existe, crear el usuario con ID 1
		superAdminUser = models.Usuario{
			ID:             1,
			PrimerNombre:   "Cristian",
			PrimerApellido: "Araya",
			EmpresaID:      empresa.ID,
			Credencial: &models.Credencial{
				ID:       1,
				Email:    "7.cristian.u@gmail.com",
				Password: "$2a$10$BryoL.Sq0BBN1efeWAQtAubPDkt.p9DTChUPc9WFOnil1.mhaZwyC", // Contraseña hasheada
			},
		}

		if result := configs.DB.Create(&superAdminUser); result.Error != nil {
			log.Fatalf("Error al crear el usuario super_administrador: %v", result.Error)
		}
	}

	// Verificar si las credenciales con IDs 1 y 2 ya están creadas
	for id := uint(1); id <= 2; id++ {
		var credencial models.Credencial
		if err := configs.DB.First(&credencial, id).Error; err != nil {
			// Si no existe, saltar (ya que están asociadas al usuario)
			log.Printf("Las credenciales necesarias con ID %d no se encontraron pero deberían estar asociadas al usuario.\n", id)
		}
	}

	// Verificar si el usuario tiene asignado el rol con ID 1
	var rolUsuario models.Rol_usuario
	if err := configs.DB.First(&rolUsuario, 1).Error; err != nil {
		// Si no existe, asignar el rol
		rolUsuario = models.Rol_usuario{
			ID:        1,
			UsuarioID: superAdminUser.ID,
			RolID:     1, // ID del rol 'super_administrador'
		}
		if result := configs.DB.Create(&rolUsuario); result.Error != nil {
			log.Fatalf("Error al asignar el rol super_administrador al usuario: %v", result.Error)
		}
	}

	log.Println("Datos iniciales creados/verificados exitosamente")
}

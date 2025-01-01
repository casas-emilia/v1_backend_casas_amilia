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
	// Verificar si el rol 'super_administrador' ya existe
	var superAdminRole models.Rol
	if err := configs.DB.Where("nombre_rol = ?", "super_administrador").First(&superAdminRole).Error; err != nil {
		// Si no existe, crear el rol
		superAdminRole = models.Rol{
			NombreRol:      "super_administrador",
			DescripcionRol: "Super administrador del sistema",
		}
		if result := configs.DB.Create(&superAdminRole); result.Error != nil {
			log.Fatalf("Error al crear el rol super_administrador: %v", result.Error)
		}
	}

	// Verificar si el rol 'administrador' ya existe
	var adminRole models.Rol
	if err := configs.DB.Where("nombre_rol = ?", "administrador").First(&adminRole).Error; err != nil {
		// Si no existe, crear el rol
		adminRole = models.Rol{
			NombreRol:      "administrador",
			DescripcionRol: "Administrador del sistema con permisos avanzados",
		}
		if result := configs.DB.Create(&adminRole); result.Error != nil {
			log.Fatalf("Error al crear el rol administrador: %v", result.Error)
		}
	}

	// Verificar si el rol 'ejecutivo_ventas' ya existe
	var salesExecutiveRole models.Rol
	if err := configs.DB.Where("nombre_rol = ?", "ejecutivo_ventas").First(&salesExecutiveRole).Error; err != nil {
		// Si no existe, crear el rol
		salesExecutiveRole = models.Rol{
			NombreRol:      "ejecutivo_ventas",
			DescripcionRol: "Responsable de gestionar las ventas",
		}
		if result := configs.DB.Create(&salesExecutiveRole); result.Error != nil {
			log.Fatalf("Error al crear el rol ejecutivo_ventas: %v", result.Error)
		}
	}

	log.Println("Roles iniciales creados correctamente")

	// Verificar si la empresa ya existe
	var empresa models.Empresa
	if err := configs.DB.Where("nombre_empresa = ?", "Empresa Demo").First(&empresa).Error; err != nil {
		// Si no existe, crear la empresa
		empresa = models.Empresa{
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

	// Verificar si el usuario 'super_administrador' ya existe
	var superAdminUser models.Usuario
	if err := configs.DB.Where("credencial->>'email' = ?", "7.cristian.u@gmail.com").First(&superAdminUser).Error; err != nil {
		// Si no existe, crear el usuario
		superAdminUser = models.Usuario{
			PrimerNombre:   "Admin",
			PrimerApellido: "Demo",
			EmpresaID:      empresa.ID,
			Credencial: &models.Credencial{
				Email:    "7.cristian.u@gmail.com",
				Password: "$10$.Dqmp/b5tmezLXRRmSDbv.W4q0.6HFeaz/.UFvLf1Qm79.zsKz9bC", // Contraseña hasheada
			},
		}
		if result := configs.DB.Create(&superAdminUser); result.Error != nil {
			log.Fatalf("Error al crear el usuario super_administrador: %v", result.Error)
		}
	}

	// Verificar si el rol ya está asignado al usuario
	var rolUsuario models.Rol_usuario
	if err := configs.DB.Where("usuario_id = ? AND rol_id = ?", superAdminUser.ID, superAdminRole.ID).First(&rolUsuario).Error; err != nil {
		// Si no está asignado, asignarlo
		rolUsuario = models.Rol_usuario{
			UsuarioID: superAdminUser.ID,
			RolID:     superAdminRole.ID,
		}
		if result := configs.DB.Create(&rolUsuario); result.Error != nil {
			log.Fatalf("Error al asignar el rol super_administrador al usuario: %v", result.Error)
		}
	}

	log.Println("Datos iniciales creados exitosamente")
}

package routers

import (
	"time"
	"v1_prefabricadas/controllers"
	"v1_prefabricadas/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Configuración de CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
			"http://192.168.0.11:3000",
			"https://v1backendcasasamilia-production.up.railway.app",
			"https://vifrontendcasasemilia-production.up.railway.app",
		}, // Dominios permitidos
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Métodos HTTP permitidos
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Headers permitidos
		ExposeHeaders:    []string{"Content-Length", "Authorization"},         // Headers expuestos al frontend
		AllowCredentials: true,                                                // Permitir cookies o credenciales
		MaxAge:           24 * time.Hour,                                      // Tiempo de caché para preflight
	}))

	/* router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}) */

	// Ruta para el login y recuperador de password(email y password :json)
	router.POST("/login", controllers.Login)
	router.POST("/password-recovery", controllers.SolicitarRecuperacion)
	router.POST("/reset-password", controllers.CambiarContrasena)

	// Rutas para tipos de estructuras
	tipos := router.Group("/tipos")
	{
		tipos.GET("/", controllers.ObtenerTipos)   // Obtener todos los Tipos de estructuras
		tipos.GET("/:id", controllers.ObtenerTipo) // Obtener Tipo estructura de acuerdo a su ID
	}

	// Rutas para categorias
	categorias := router.Group("/categorias")
	{
		categorias.GET("/", controllers.ObtenerCategorias)   // Obtener todas las Categorias
		categorias.GET("/:id", controllers.ObtenerCategoria) // Obtener Categoria de acuerdo al ID
	}

	estilos := router.Group("/estilos")
	{
		estilos.GET("/", controllers.ObtenerEstilos)   // Obtener todos los Estilos
		estilos.GET("/:id", controllers.ObtenerEstilo) // Obtener Estilo de acuerdo a su ID
	}

	empresas := router.Group("/empresas")
	{
		empresas.GET("/", controllers.ObtenerEmpresas)          // Obtener todas las Empresas
		empresas.GET("/:empresaID", controllers.ObtenerEmpresa) // Obtener datos de Empresa de acuerdo a su ID

		servicios := empresas.Group("/:empresaID/servicios")
		{
			servicios.GET("/", controllers.ObtenerServicios)           // Obtener todos los Servicios
			servicios.GET("/:servicioID", controllers.ObtenerServicio) // Obtener Servicio de acuerdo a su ID
		}

		redes := empresas.Group("/:empresaID/redes")
		{
			redes.GET("/", controllers.ObtenerRedes)     // Obtener todas las redes sociales de la empresa
			redes.GET("/:redID", controllers.ObtenerRed) // Obtener Red social de acuerdo a su ID
		}

		portadas := empresas.Group("/:empresaID/portadas")
		{
			portadas.GET("/", controllers.ObtenerPortadas)          // Obtener todas las portadas de la Empresa
			portadas.GET("/:portadaID", controllers.ObtenerPortada) // Obtener Portada de acuerdo a su ID
		}

		noticiasEmpresa := empresas.Group("/:empresaID/noticiasEmpresa")
		{
			noticiasEmpresa.GET("/", controllers.ObtenerNoticiasEmpresa)          // Función para obtener todas las noticias de una empresa
			noticiasEmpresa.GET("/:noticiaID", controllers.ObtenerNoticiaEmpresa) // Función para obtener una noticia de empresa

			imagenesNoticiasEmpresa := noticiasEmpresa.Group("/:noticiaID/imagenesNoticiasEmpresa")
			{
				imagenesNoticiasEmpresa.GET("/", controllers.ObtenerImagenesNoticias)              // Obtener todas las imagenes de una Noticia
				imagenesNoticiasEmpresa.GET("/:imagenNoticiaID", controllers.ObtenerImagenNoticia) // Obtener una Imagen de una Noticia
			}
		}

		prefabricadas := empresas.Group("/:empresaID/prefabricadas")
		{
			prefabricadas.GET("/", controllers.ObtenerPrefabricadas)               // Obtener todas las Prefabricadas de la Empresa
			prefabricadas.GET("", controllers.ObtenerPrefabricadas)                // Obtener todas las Prefabricadas de la Empresa (sin slash al final)
			prefabricadas.GET("/:prefabricadaID", controllers.ObtenerPrefabricada) // Obtener Prefabricada de acuerdo a su ID

			imagenesPrefabricadas := prefabricadas.Group("/:prefabricadaID/imagenesPrefabricadas")
			{
				imagenesPrefabricadas.GET("/", controllers.ObtenerImagenesPrefabricadas)                  // Obtener todas las Imagenes de una Prefabricada
				imagenesPrefabricadas.GET("/:imagenPrefabricadaID", controllers.ObtenerImagePrefabricada) // Obtener una imagen de acuerdo a su ID de la prefabricada
			}

			caracteristicas := prefabricadas.Group("/:prefabricadaID/caracteristicas")
			{
				caracteristicas.GET("/", controllers.ObtenerCaracteristicas)                 // Obtener todas las características de la Prefabricada
				caracteristicas.GET("/:caracteristicaID", controllers.ObtenerCaracteristica) // Obtener característica de acuerdo al ID

			}

			precios := prefabricadas.Group("/:prefabricadaID/precios")
			{
				precios.GET("/", controllers.ObtenerPrecios)         // Otener todos los Precios de una Prefabricada
				precios.GET("/:precioID", controllers.ObtenerPrecio) // Obtener un Precio de acuerdo a su ID

				incluyes := precios.Group("/:precioID/incluyes")
				{
					incluyes.GET("/", controllers.ObtenerIncluyes)          // Obtener todos los Incluyes de un Precio
					incluyes.GET("/:incluyeID", controllers.ObtenerIncluye) // Obtener Incluye de acuerdo a su ID
				}
			}
		}

		usuarios := empresas.Group("/:empresaID/usuarios")
		{
			usuarios.GET("/", controllers.ObtenerUsuarios)          // Obtener todos los Usuarios de una Empresa
			usuarios.GET("/:usuarioID", controllers.ObtenerUsuario) // Obtener Usuario de acuerdo a su ID

			contactos := usuarios.Group("/:usuarioID/contactos")
			{
				contactos.GET("/", controllers.ObtenerContactos)           // Obtener todos los Contactos de un Usuario
				contactos.GET("/:contactoID", controllers.ObtenerContacto) // Obtener Contacto de acuerdo a su ID
			}

		}
	}

	roles := router.Group("/roles")
	{
		roles.GET("/", controllers.ObtenerRoles)     // Obtener todos los Roles
		roles.GET("/:rolID", controllers.ObtenerRol) // Obtener Rol de acuerdo a su ID

		rolesUsuarios := roles.Group("/:rolID/roles_usuarios")
		{
			rolesUsuarios.GET("/", controllers.ObtenerRoles_usuarios)            // Obtener todos los roles y sus usuarios
			rolesUsuarios.GET("/:rol_usuarioID", controllers.ObtenerRol_usuario) // Obtener Rol_usuario
		}
	}

	// Rutas Administración del sistema
	admin := router.Group("/administracion", middlewares.AuthMiddleware(), middlewares.SuperAdminMiddleware())
	{
		admin.GET("/", controllers.AdminObtenerServicios) // Página principal del panel

		// Rutas para tipos de estructuras
		tipos := admin.Group("/tipos")
		{
			tipos.POST("/", controllers.CrearTipo)         // Crear un Tipo de estructura
			tipos.GET("/", controllers.ObtenerTipos)       // Obtener todos los Tipos de estructuras
			tipos.GET("/:id", controllers.ObtenerTipo)     // Obtener Tipo estructura de acuerdo a su ID
			tipos.PUT("/:id", controllers.ActualizarTipo)  // Actualizar datos de Tipos de estructuras
			tipos.DELETE("/:id", controllers.EliminarTipo) // Eliminar Tipo estructura de acuerdo a su ID(Eliminación lógica)
		}

		// Rutas para categorias
		categorias := admin.Group("/categorias")
		{
			categorias.POST("/", controllers.CrearCategoria)         // Crear Categoria
			categorias.GET("/", controllers.ObtenerCategorias)       // Obtener todas las Categorias
			categorias.GET("/:id", controllers.ObtenerCategoria)     // Obtener Categoria de acuerdo al ID
			categorias.PUT("/:id", controllers.ActualizarCategoria)  // Actualizar datos de Categoría
			categorias.DELETE("/:id", controllers.EliminarCategoria) // Eliminar lógicamente Categoria de acuerdo a su ID
		}

		estilos := admin.Group("/estilos")
		{
			estilos.POST("/", controllers.CrearEstilo)         // Crear Estilo
			estilos.GET("/", controllers.ObtenerEstilos)       // Obtener todos los Estilos
			estilos.GET("/:id", controllers.ObtenerEstilo)     // Obtener Estilo de acuerdo a su ID
			estilos.PUT("/:id", controllers.ActualizarEstilo)  // Actualizar datos de Estilo de acuerdo a su ID
			estilos.DELETE("/:id", controllers.EliminarEstilo) // Eliminar lógicamente un Estilo de acuerdo al ID
		}

		empresas := admin.Group("/empresas")
		{
			empresas.POST("/", controllers.CrearEmpresa)                // Crear Empresa
			empresas.GET("/", controllers.ObtenerEmpresas)              // Obtener todas las Empresas
			empresas.GET("/:empresaID", controllers.ObtenerEmpresa)     // Obtener datos de Empresa de acuerdo a su ID
			empresas.PUT("/:empresaID", controllers.ActualizarEmpresa)  // Actualizar datos de Empresa
			empresas.DELETE("/:empresaID", controllers.EliminarEmpresa) // Eliminar Empresa de acuerdo a su ID

			servicios := empresas.Group("/:empresaID/servicios")
			{
				servicios.POST("/", controllers.CrearServicio)                 // Crear un servicio de la Empresa
				servicios.GET("/", controllers.ObtenerServicios)               // Obtener todos los Servicios
				servicios.GET("/:servicioID", controllers.ObtenerServicio)     // Obtener Servicio de acuerdo a su ID
				servicios.PUT("/:servicioID", controllers.ActualizarServicio)  // Actualizar datos del Servicio de acuerdo a su ID
				servicios.DELETE("/:servicioID", controllers.EliminarServicio) // Eliminar lógicamente un servicio de acuerdo a su ID
			}

			redes := empresas.Group("/:empresaID/redes")
			{
				redes.POST("/", controllers.CrearRed)            // Crea Red Social de la Empresa
				redes.GET("/", controllers.ObtenerRedes)         // Obtener todas las redes sociales de la empresa
				redes.GET("/:redID", controllers.ObtenerRed)     // Obtener Red social de acuerdo a su ID
				redes.PUT("/:redID", controllers.ActualizarRed)  // Actualizar datos de una red social de acuerdo a su ID
				redes.DELETE("/:redID", controllers.EliminarRed) // Eliminar Red Social de acuerdo al ID enviado
			}

			portadas := empresas.Group("/:empresaID/portadas")
			{
				portadas.POST("/", controllers.CrearPortada)                // Crear Portada
				portadas.GET("/", controllers.ObtenerPortadas)              // Obtener todas las portadas de la Empresa
				portadas.GET("/:portadaID", controllers.ObtenerPortada)     // Obtener Portada de acuerdo a su ID
				portadas.PUT("/:portadaID", controllers.ActualizarPortada)  // Actualizar datos de una Portada
				portadas.DELETE("/:portadaID", controllers.EliminarPortada) // Eliminar una portada
			}

			noticiasEmpresa := empresas.Group("/:empresaID/noticiasEmpresa")
			{
				noticiasEmpresa.POST("/", controllers.CrearNoticia)                   // Crear una Noticia
				noticiasEmpresa.GET("/", controllers.ObtenerNoticiasEmpresa)          // Función para obtener todas las noticias de una empresa
				noticiasEmpresa.GET("/:noticiaID", controllers.ObtenerNoticiaEmpresa) // Función para obtener una noticia de empresa
				noticiasEmpresa.PUT("/:noticiaID", controllers.ActualizarNoticia)     // Actualizar datos de Una Noticias
				noticiasEmpresa.DELETE("/:noticiaID", controllers.EliminarNoticia)    // Eliminar lógicamente una Noticia de acuerdo a su ID

				imagenesNoticiasEmpresa := noticiasEmpresa.Group("/:noticiaID/imagenesNoticiasEmpresa")
				{
					imagenesNoticiasEmpresa.POST("/", controllers.CrearImagenNoticia)                      // Crear una imagen para una Noticia
					imagenesNoticiasEmpresa.GET("/", controllers.ObtenerImagenesNoticias)                  // Obtener todas las imagenes de una Noticia
					imagenesNoticiasEmpresa.GET("/:imagenNoticiaID", controllers.ObtenerImagenNoticia)     // Obtener una Imagen de una Noticia
					imagenesNoticiasEmpresa.PUT("/:imagenNoticiaID", controllers.ActualizarImagenNoticia)  // Actualizar una imagen de una Noticia
					imagenesNoticiasEmpresa.DELETE("/:imagenNoticiaID", controllers.EliminarImagenNoticia) // Eliminar Logicamente una Imagen de una Noticia
				}
			}

			prefabricadas := empresas.Group("/:empresaID/prefabricadas")
			{
				prefabricadas.POST("/", controllers.CrearPrefabricada)                     // Crear una prefabricada
				prefabricadas.GET("/", controllers.ObtenerPrefabricadas)                   // Obtener todas las Prefabricadas de la Empresa
				prefabricadas.GET("", controllers.ObtenerPrefabricadas)                    // Obtener todas las Prefabricadas de la Empresa (sin slash al final)
				prefabricadas.GET("/:prefabricadaID", controllers.ObtenerPrefabricada)     // Obtener Prefabricada de acuerdo a su ID
				prefabricadas.PUT("/:prefabricadaID", controllers.ActualizarPrefabricada)  // Actualizar datos de Prefabricada
				prefabricadas.DELETE("/:prefabricadaID", controllers.EliminarPrefabricada) // Eliminar lógicamente una Prefabricada de acuerdo al ID enviado

				imagenesPrefabricadas := prefabricadas.Group("/:prefabricadaID/imagenesPrefabricadas")
				{
					imagenesPrefabricadas.POST("/", controllers.CrearImagen_prefabricada)                          // Crear imagen prefabricada
					imagenesPrefabricadas.GET("/", controllers.ObtenerImagenesPrefabricadas)                       // Obtener todas las Imagenes de una Prefabricada
					imagenesPrefabricadas.GET("/:imagenPrefabricadaID", controllers.ObtenerImagePrefabricada)      // Obtener una imagen de acuerdo a su ID de la prefabricada
					imagenesPrefabricadas.PUT("/:imagenPrefabricadaID", controllers.ActualizarImagenPrefabricada)  // Actualizar datos de una imagen
					imagenesPrefabricadas.DELETE("/:imagenPrefabricadaID", controllers.EliminarImagenPrefabricada) // Eliminar lógicamente una imagen_prefabricada
				}

				caracteristicas := prefabricadas.Group("/:prefabricadaID/caracteristicas")
				{
					caracteristicas.POST("/", controllers.CrearCaracteristica)                       // Crear una nueva característica
					caracteristicas.GET("/", controllers.ObtenerCaracteristicas)                     // Obtener todas las características de la Prefabricada
					caracteristicas.GET("/:caracteristicaID", controllers.ObtenerCaracteristica)     // Obtener característica de acuerdo al ID
					caracteristicas.PUT("/:caracteristicaID", controllers.ActualizarCaracteristica)  // Actualizar datos de Caracterítica
					caracteristicas.DELETE("/:caracteristicaID", controllers.EliminarCaracteristica) // Eliminar lógicamente una Característica
				}

				precios := prefabricadas.Group("/:prefabricadaID/precios")
				{
					precios.POST("/", controllers.CrearPrecio)               // Crear un precio a una Prefabricada
					precios.GET("/", controllers.ObtenerPrecios)             // Otener todos los Precios de una Prefabricada
					precios.GET("/:precioID", controllers.ObtenerPrecio)     // Obtener un Precio de acuerdo a su ID
					precios.PUT("/:precioID", controllers.ActualizarPrecio)  // Actualizar un Precio de acuerdo al ID
					precios.DELETE("/:precioID", controllers.EliminarPrecio) // Eliminar un Precio de acuerdo al ID

					incluyes := precios.Group("/:precioID/incluyes")
					{
						incluyes.POST("/", controllers.CrearIncluye)                // Crear un Incluye de un precio
						incluyes.GET("/", controllers.ObtenerIncluyes)              // Obtener todos los Incluyes de un Precio
						incluyes.GET("/:incluyeID", controllers.ObtenerIncluye)     // Obtener Incluye de acuerdo a su ID
						incluyes.PUT("/:incluyeID", controllers.ActualizarIncluye)  // Actualizar Incluye
						incluyes.DELETE("/:incluyeID", controllers.EliminarIncluye) // Eliminar lógicamente un Incluye
					}
				}
			}

			usuarios := empresas.Group("/:empresaID/usuarios")
			{
				usuarios.POST("/", controllers.CrearUsuario)                // Crear un nuevo usuarios
				usuarios.GET("/", controllers.ObtenerUsuarios)              // Obtener todos los Usuarios de una Empresa
				usuarios.GET("/:usuarioID", controllers.ObtenerUsuario)     // Obtener Usuario de acuerdo a su ID
				usuarios.PUT("/:usuarioID", controllers.ActualizarUsuario)  // Actualizar datos de Usuario de acuerdo a su ID
				usuarios.DELETE("/:usuarioID", controllers.EliminarUsuario) // Eliminar Usuario lógicamente de acuerdo a su ID

				contactos := usuarios.Group("/:usuarioID/contactos")
				{
					contactos.POST("/", controllers.CrearContacto)                 // Crear datos de Contacto de Usuario
					contactos.GET("/", controllers.ObtenerContactos)               // Obtener todos los Contactos de un Usuario
					contactos.GET("/:contactoID", controllers.ObtenerContacto)     // Obtener Contacto de acuerdo a su ID
					contactos.PUT("/:contactoID", controllers.ActualizarContacto)  // Actualizar datos de Contacto
					contactos.DELETE("/:contactoID", controllers.EliminarContacto) // Eliminar Datos de Contacto lógicamente
				}

				credenciales := usuarios.Group("/:usuarioID/credenciales")
				{
					credenciales.GET("/", controllers.ObtenerCredenciales)                 // Obtener credenciales de un usuario
					credenciales.POST("/", controllers.CrearCredencial)                    // Crear Credenciales de Usuario
					credenciales.PUT("/:credencialID", controllers.ActualizarCredenciales) // Actualizar credenciales de acceso
				}

			}
		}

		roles := admin.Group("/roles")
		{
			roles.POST("/", controllers.CrearRol)            // Crea un nuevo Rol
			roles.GET("/", controllers.ObtenerRoles)         // Obtener todos los Roles
			roles.GET("/:rolID", controllers.ObtenerRol)     // Obtener Rol de acuerdo a su ID
			roles.PUT("/:rolID", controllers.ActualizarRol)  // Actualizar datos de un Rol
			roles.DELETE("/:rolID", controllers.EliminarRol) // Eliminar Rol lógicamente

			rolesUsuarios := roles.Group("/:rolID/roles_usuarios")
			{
				rolesUsuarios.POST("/", controllers.CrearRol_usuario)                    // Asignarle un Rol a un Usuario
				rolesUsuarios.GET("/", controllers.ObtenerRoles_usuarios)                // Obtener todos los roles y sus usuarios
				rolesUsuarios.GET("/:rol_usuarioID", controllers.ObtenerRol_usuario)     // Obtener Rol_usuario
				rolesUsuarios.PUT("/:rol_usuarioID", controllers.ActualizarRol_usuario)  // Actualizar datos Rol_usuario
				rolesUsuarios.DELETE("/:rol_usuarioID", controllers.EliminarRol_usuario) // Eliminar lógicamente un rol_usuario
			}
		}
	}

	return router
}

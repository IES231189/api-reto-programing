package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"persona_api/src/Persona/infraestructure/handler"
	"persona_api/src/Persona/infraestructure/repository"
	"persona_api/src/Persona/aplication"
	"persona_api/src/Persona/domain"
)

func cargarVariablesEntorno() {
	if err := godotenv.Load(); err != nil {
		log.Println(" No se encontró un archivo .env o hubo un error al cargarlo")
	}
}

func conectarBD() *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar la base de datos:", err)
	}

	db.AutoMigrate(&domain.Persona{})
	return db
}

func main() {
	cargarVariablesEntorno()
	db := conectarBD()

	// Inyección de dependencias
	repo := repository.NewPersonaRepository(db)
	service := application.NewPersonaService(repo)
	handler := handler.NewPersonaHandler(service)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, 
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, 
	}))

	r.POST("/personas", handler.CrearPersona)
	r.GET("/personas", handler.ObtenerPersonas)
	r.GET("/personas/conteo-genero", handler.ContarGeneros)
	r.GET("/personas/conteo-genero-longpoll", handler.ContarGenerosLongPolling)

	puerto := os.Getenv("SERVER_PORT")
	if puerto == "" {
		puerto = "4040" 
	}

	fmt.Println("Servidor corriendo en http://localhost:" + puerto)
	r.Run(":" + puerto)
}
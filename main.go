package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"persona_api/src/Persona/infraestructure/handler"
	"persona_api/src/Persona/infraestructure/repository"
	"persona_api/src/Persona/aplication"
	"persona_api/src/Persona/domain"
)

func conectarBD() *gorm.DB {
	dsn := "appuser:SuperClave@tcp(127.0.0.1:3306)/hexagonal_db?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar la base de datos:", err)
	}

	db.AutoMigrate(&domain.Persona{})

	return db
}

func main() {
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
	r.GET("/personas/conteo-genero-longpoll", handler.ContarGenerosLongPolling)
	
	fmt.Println("Servidor corriendo en http://localhost:8080")
	r.Run(":8080")
}
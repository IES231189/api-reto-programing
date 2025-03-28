// domain/persona.go
package domain

type Persona struct {
	ID     uint
	Nombre string
	Edad   int
	Sexo   string `gorm:"size:10"`  // Masculino/Femenino
	Genero string `gorm:"size:15"`   // Identidad de género
}
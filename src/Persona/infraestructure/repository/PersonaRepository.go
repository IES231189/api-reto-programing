
package repository

import (
	"gorm.io/gorm"
	"persona_api/src/Persona/domain"
)

type PersonaRepository interface {
	Guardar(persona *domain.Persona) error
	ObtenerTodas() ([]domain.Persona, error)
	ContarPorGenero() (map[string]int, error)
}

type personaRepository struct {
	db *gorm.DB
}

func NewPersonaRepository(db *gorm.DB) PersonaRepository {
	return &personaRepository{db}
}

func (r *personaRepository) Guardar(persona *domain.Persona) error {
	return r.db.Create(persona).Error
}

func (r *personaRepository) ObtenerTodas() ([]domain.Persona, error) {
	var personas []domain.Persona
	err := r.db.Find(&personas).Error
	return personas, err
}

func (r *personaRepository) ContarPorGenero() (map[string]int, error) {
	var resultados []struct {
		Genero string
		Count  int
	}
	
	err := r.db.Model(&domain.Persona{}).
		Select("genero, count(*) as count").
		Group("genero").
		Scan(&resultados).Error
	
	if err != nil {
		return nil, err
	}
	
	conteo := make(map[string]int)
	for _, res := range resultados {
		conteo[res.Genero] = res.Count
	}
	
	return conteo, nil
}
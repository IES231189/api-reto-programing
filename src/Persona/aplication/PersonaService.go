package application

import (
	"persona_api/src/Persona/domain"
	"persona_api/src/Persona/infraestructure/repository"
	"time"
)

type PersonaService interface {
	CrearPersona(nombre string, edad int, sexo string, genero string) (*domain.Persona, error)
	ListarPersonas() ([]domain.Persona, error)
	ContarPorGenero() (map[string]int, error)
	ContarPorGeneroLongPolling(timeout time.Duration, cambios chan<- map[string]int) error
}

type personaService struct {
	repo repository.PersonaRepository
}

func NewPersonaService(repo repository.PersonaRepository) PersonaService {
	return &personaService{repo}
}

func (s *personaService) CrearPersona(nombre string, edad int, sexo string, genero string) (*domain.Persona, error) {
	persona := &domain.Persona{
		Nombre: nombre,
		Edad:   edad,
		Sexo:   sexo,
		Genero: genero,
	}
	err := s.repo.Guardar(persona)
	return persona, err
}

func (s *personaService) ListarPersonas() ([]domain.Persona, error) {
	return s.repo.ObtenerTodas()
}

func (s *personaService) ContarPorGenero() (map[string]int, error) {
	return s.repo.ContarPorGenero()
}

func (s *personaService) ContarPorGeneroLongPolling(timeout time.Duration, cambios chan<- map[string]int) error {
	ultimoConteo := make(map[string]int)
	
	for {
		select {
		case <-time.After(timeout):
			return nil // Timeout sin cambios
		default:
			conteoActual, err := s.repo.ContarPorGenero()
			if err != nil {
				return err
			}
			
			// Verificar si hay cambios
			if !mapasIguales(ultimoConteo, conteoActual) {
				cambios <- conteoActual
				ultimoConteo = conteoActual
				return nil // Cambio detectado, terminamos
			}
			
			time.Sleep(1 * time.Second) // Esperar antes de revisar de nuevo
		}
	}
}

func mapasIguales(a, b map[string]int) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
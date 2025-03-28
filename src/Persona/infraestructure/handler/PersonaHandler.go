
package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"persona_api/src/Persona/aplication"
)

type PersonaHandler struct {
	service application.PersonaService
}

func NewPersonaHandler(service application.PersonaService) *PersonaHandler {
	return &PersonaHandler{service}
}

func (h *PersonaHandler) CrearPersona(c *gin.Context) {
	var request struct {
		Nombre string `json:"nombre"`
		Edad   int    `json:"edad"`
		Sexo   string `json:"sexo"`
		Genero string `json:"genero"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	persona, err := h.service.CrearPersona(request.Nombre, request.Edad, request.Sexo, request.Genero)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo crear la persona"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Agregado correctamente",
		"contenido": persona,
	})
}

func (h *PersonaHandler) ObtenerPersonas(c *gin.Context) {
	personas, err := h.service.ListarPersonas()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pueden obtener personas"})
		return
	}

	c.JSON(http.StatusOK, personas)
}

func (h *PersonaHandler) ContarGeneros(c *gin.Context) {
	conteo, err := h.service.ContarPorGenero()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No se pudo contar por género"})
		return
	}

	c.JSON(http.StatusOK, conteo)
}

func (h *PersonaHandler) ContarGenerosLongPolling(c *gin.Context) {
	timeoutStr := c.DefaultQuery("timeout", "30")
	timeout, err := time.ParseDuration(timeoutStr + "s")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Timeout inválido"})
		return
	}

	cambios := make(chan map[string]int)
	
	go func() {
		err := h.service.ContarPorGeneroLongPolling(timeout, cambios)
		if err != nil {
			close(cambios)
		}
	}()

	select {
	case conteo, ok := <-cambios:
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error en el servidor"})
			return
		}
		c.JSON(http.StatusOK, conteo)
	case <-time.After(timeout):
		c.JSON(http.StatusRequestTimeout, gin.H{"message": "No hubo cambios en el período especificado"})
	}
}
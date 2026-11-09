package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/domain"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/dto"
	"github.com/ArtemOlyinik/client_server_Practichna_05/internal/middleware"
)

type IncidentHandler struct {
	db *gorm.DB
}

func NewIncidentHandler(db *gorm.DB) *IncidentHandler {
	return &IncidentHandler{db: db}
}

// List godoc
// @Summary Отримати всі інциденти
// @Tags incidents
// @Produce json
// @Success 200 {array} domain.Incident
// @Failure 500 {object} dto.ErrorResponse
// @Router /incidents [get]
func (h *IncidentHandler) List(c *gin.Context) {
	var incidents []domain.Incident
	if err := h.db.Order("id asc").Find(&incidents).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}
	c.JSON(http.StatusOK, incidents)
}

// Create godoc
// @Summary Створити інцидент
// @Tags incidents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateIncidentDTO true "Дані інциденту"
// @Success 201 {object} domain.Incident
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /incidents [post]
func (h *IncidentHandler) Create(c *gin.Context) {
	var input dto.CreateIncidentDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation Error",
			Details: validationErrors(err),
		})
		return
	}

	incident := domain.Incident{
		Title:       input.Title,
		Description: input.Description,
		Priority:    input.Priority,
		Status:      "open",
		UserID:      middleware.GetUserID(c),
	}

	if err := h.db.Create(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusCreated, incident)
}

// Get godoc
// @Summary Отримати інцидент за ID
// @Tags incidents
// @Produce json
// @Param id path int true "ID інциденту"
// @Success 200 {object} domain.Incident
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /incidents/{id} [get]
func (h *IncidentHandler) Get(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid incident ID"})
		return
	}

	var incident domain.Incident
	if err := h.db.First(&incident, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Incident not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// UpdateStatus godoc
// @Summary Змінити статус інциденту
// @Tags incidents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID інциденту"
// @Param request body dto.UpdateStatusDTO true "Новий статус"
// @Success 200 {object} domain.Incident
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /incidents/{id}/status [patch]
func (h *IncidentHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid incident ID"})
		return
	}

	var input dto.UpdateStatusDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Validation Error",
			Details: validationErrors(err),
		})
		return
	}

	var incident domain.Incident
	if err := h.db.First(&incident, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Incident not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	incident.Status = input.Status
	if err := h.db.Save(&incident).Error; err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}

	c.JSON(http.StatusOK, incident)
}

// Delete godoc
// @Summary Видалити інцидент
// @Tags incidents
// @Security BearerAuth
// @Param id path int true "ID інциденту"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /incidents/{id} [delete]
func (h *IncidentHandler) Delete(c *gin.Context) {
	id, err := parseID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: "Invalid incident ID"})
		return
	}

	result := h.db.Delete(&domain.Incident{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: "Database error"})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{Error: "Incident not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

func parseID(value string) (uint64, error) {
	return strconv.ParseUint(value, 10, 64)
}

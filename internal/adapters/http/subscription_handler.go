package http

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tamaqazaq/subscription-service/internal/domain"
	"github.com/tamaqazaq/subscription-service/internal/domain/application"
	"log"
	"net/http"
	"time"
)

type SubscriptionHandler struct {
	usecase application.SubscriptionUsecase
}

func NewSubscriptionHandler(u application.SubscriptionUsecase) *SubscriptionHandler {
	return &SubscriptionHandler{usecase: u}
}

// @Summary Create subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body domain.Subscription true "Subscription data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var sub domain.Subscription

	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if sub.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	if sub.StartDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing or invalid start_date"})
		return
	}

	if err := h.usecase.Create(&sub); err != nil {
		log.Println("CreateSubscription error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": sub.ID.String()})
}

// @Summary Get all subscriptions
// @Tags Subscriptions
// @Produce json
// @Success 200 {array} domain.Subscription
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) GetAll(c *gin.Context) {
	subs, err := h.usecase.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subs)
}

// @Summary Get subscription by ID
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} domain.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	sub, err := h.usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

// @Summary Update subscription
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body domain.Subscription true "Updated data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}

	var sub domain.Subscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	if err := h.usecase.Update(id, &sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

// @Summary Delete subscription
// @Tags Subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid UUID"})
		return
	}
	if err := h.usecase.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// GetTotal handles the total cost of subscriptions within a period.
//
// @Summary Get total subscription cost
// @Description Calculates the sum of all subscriptions in a given period. Optional filters: user_id, service_name.
// @Tags Subscriptions
// @Produce json
// @Param user_id query string false "User ID (optional UUID)"
// @Param service_name query string false "Service name (optional)"
// @Param start query string true "Start date in format MM-YYYY (e.g. 07-2025)"
// @Param end query string true "End date in format MM-YYYY (e.g. 12-2025)"
// @Success 200 {object} map[string]int "e.g. {\"total\": 1200}"
// @Failure 400 {object} map[string]string "e.g. {\"error\": \"Invalid start format (MM-YYYY)\"}"
// @Failure 500 {object} map[string]string "e.g. {\"error\": \"Failed to calculate total\"}"
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) GetTotal(c *gin.Context) {
	userIDStr := c.Query("user_id")
	serviceName := c.Query("service_name")
	startStr := c.Query("start")
	endStr := c.Query("end")

	var userID *uuid.UUID
	if userIDStr != "" {
		parsed, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
		userID = &parsed
	}

	var servicePtr *string
	if serviceName != "" {
		servicePtr = &serviceName
	}

	start, err := time.Parse("01-2006", startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start format (MM-YYYY)"})
		return
	}
	end, err := time.Parse("01-2006", endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end format (MM-YYYY)"})
		return
	}

	total, err := h.usecase.GetTotal(userID, servicePtr, start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}

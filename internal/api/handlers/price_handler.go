package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PriceHandler struct {
	PriceHandler *service.PriceService
}

func NewPriceHandler(priceHandler *service.PriceService) *PriceHandler {
	return &PriceHandler{PriceHandler: priceHandler}
}

func (ph *PriceHandler) CreatePriceForShowtime(c *gin.Context) {
	var req model.Price
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := ph.PriceHandler.CreatePriceForShowtime(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation successful"})
}

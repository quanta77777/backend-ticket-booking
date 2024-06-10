package handlers

import (
	"net/http"
	"strconv"

	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type TicketHandler struct {
	ticketService *service.TicketService
}

func NewTicketHandler(ticketService *service.TicketService) *TicketHandler {
	return &TicketHandler{ticketService: ticketService}
}

func (h *TicketHandler) CreateTicket(c *gin.Context) {
	var ticket model.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticketID, err := h.ticketService.CreateTicket(&ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"ticket_id": ticketID})
}

func (h *TicketHandler) AddProductWithTicketId(c *gin.Context) {
	var ticket *model.TicketProduct
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ticketService.AddProductWithTicketId(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error add": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Add product with ticket id  successfully")
}

func (h *TicketHandler) AddSeatWithTicketId(c *gin.Context) {
	var ticket *model.TicketSeat
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.ticketService.AddSeatWithTicketId(ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, "Add seat with ticket id  successfully")
}

func (h *TicketHandler) UserHasTicketForMovie(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	hasTicket, err := h.ticketService.UserHasTicketForMovie(userID, movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"has_ticket": hasTicket})
}

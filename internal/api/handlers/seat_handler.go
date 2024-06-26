package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	SeatHandler *service.SeatService
}

func NewSeatHandler(seatHandler *service.SeatService) *SeatHandler {
	return &SeatHandler{SeatHandler: seatHandler}
}

func (sh *SeatHandler) AddSeat(c *gin.Context) {
	var seatReq struct {
		TheaterId  int    `json:"theater_id"`
		SeatType   string `json:"seat_type"`
		SeatNumber string `json:"seat_number"`
	}

	if err := c.BindJSON(&seatReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Invalid req"})
	}
	err := sh.SeatHandler.SeatService.AddSeat(seatReq.TheaterId, seatReq.SeatType, seatReq.SeatNumber)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": "Err from server add seat"})
	}

	c.JSON(http.StatusCreated, gin.H{"status": "Add cuccessfully"})

}

func (sh *SeatHandler) GetSeatBytheaterID(c *gin.Context) {
	theaterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater ID"})
		return
	}

	seats, err := sh.SeatHandler.GetSeatByTheaterID(theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seats)
}

func (sh *SeatHandler) GetSeatsWithReservationStatusAndPrices(c *gin.Context) {
	theaterID, err := strconv.Atoi(c.Query("theater_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater_id"})
		return
	}
	showtimeID, err := strconv.Atoi(c.Query("showtime_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime_id"})
		return
	}

	seats, prices, err := sh.SeatHandler.GetSeatsWithReservationStatusAndPrices(theaterID, showtimeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	response := gin.H{
		"seats":  seats,
		"prices": prices,
	}
	c.JSON(http.StatusOK, response)
}

func (sh *SeatHandler) ReserveSeat(c *gin.Context) {
	var reservation model.SeatReservation
	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sh.SeatHandler.CreateReservation(reservation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation successful"})
}

func (sh *SeatHandler) GetSeatByTicketId(c *gin.Context) {
	ticketID, err := strconv.Atoi(c.Param("ticket_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	seats, err := sh.SeatHandler.GetSeatByTicketId(ticketID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, seats)
}

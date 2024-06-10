package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TheaterHandler struct {
	TheaterService *service.TheaterService
}

func NewTheaterHandler(theaterService *service.TheaterService) *TheaterHandler {
	return &TheaterHandler{TheaterService: theaterService}
}

func (th *TheaterHandler) AddTheater(c *gin.Context) {
	var theaterInput struct {
		BranchId int    `json:"branch_id"`
		Name     string `json:"name"`
	}

	if err := c.BindJSON(&theaterInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	TheaterID, err := th.TheaterService.AddTheater(theaterInput.BranchId, theaterInput.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add theater"})
	}

	newTheater := model.Theater{
		TheaterId: int(TheaterID),
		Name:      theaterInput.Name,
	}
	c.JSON(http.StatusCreated, newTheater)

}

func (th *TheaterHandler) GetAllTheaterByBranchID(c *gin.Context) {
	theaterIDStr := c.Param("id")
	theaterID, err := strconv.Atoi(theaterIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	theaters, err := th.TheaterService.GetTheatersByBranchID(theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, theaters)
}

func (th *TheaterHandler) GetTheaterByID(c *gin.Context) {
	theaterID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid theater ID"})
		return
	}
	theater, err := th.TheaterService.GetTheaterByID(theaterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, theater)
}

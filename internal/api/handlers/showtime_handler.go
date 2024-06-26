package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ShowtimeHandler struct {
	ShowtimeService *service.ShowtimeService
}

func NewShowtimeHandler(showtimeService *service.ShowtimeService) *ShowtimeHandler {
	return &ShowtimeHandler{ShowtimeService: showtimeService}
}

func (mh *ShowtimeHandler) GetShowtimeByDay(c *gin.Context) {
	// Retrieve and validate the day query parameter
	dayStr := c.Query("day")
	day, err := time.Parse("2006-01-02", dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	// Retrieve the movie_id query parameter
	movieIDStr := c.Query("movie_id")
	var movieID int
	if movieIDStr != "" {
		movieID, err = strconv.Atoi(movieIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie_id"})
			return
		}
	}

	var showtimes []model.Showtime

	// Retrieve showtimes based on the presence of MovieID
	showtimes, err = mh.ShowtimeService.GetShowtimeByDayAndMovieID(day, movieID)
	// if movieID > 0 {
	// 	showtimes, err = mh.ShowtimeService.GetShowtimeByDayAndMovieID(day, movieID)
	// }
	//  else {
	// 	showtimes, err = mh.ShowtimeService.GetAllShowtimeByDay(day)
	// }

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, showtimes)

}

func (mh *ShowtimeHandler) GetShowtimeWithBranch(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}
	var request struct {
		Day string `form:"day" binding:"required"`
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorr": err.Error()})
		return
	}
	day, err := time.Parse("2006-01-02", request.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	var showtimes []model.Showtime
	showtimes, err = mh.ShowtimeService.GetShowtimeWithBranch(branchID, day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)

}

func (mh *ShowtimeHandler) GetShowtimeWithCinema(c *gin.Context) {
	cinemaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}
	var request struct {
		Day string `form:"day" binding:"required"`
	}
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errorr": err.Error()})
		return
	}
	day, err := time.Parse("2006-01-02", request.Day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}
	var showtimes []model.Showtime
	showtimes, err = mh.ShowtimeService.GetShowtimeWithCinema(cinemaID, day)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)

}

func (mh *ShowtimeHandler) GetShowtimeByID(c *gin.Context) {
	showtimeID, err := strconv.Atoi(c.Param("showtime_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid showtime ID"})
		return
	}

	var showtimes []model.Showtime
	showtimes, err = mh.ShowtimeService.GetShowtimByID(showtimeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, showtimes)

}

func (mh *ShowtimeHandler) AddShowtime(c *gin.Context) {
	var ShowtimeReq struct {
		BranchID  int       `json:"branch_id"`
		MovieID   int       `json:"movie_id"`
		TheaterID int       `json:"theater_id"`
		CinemaID  int       `json:"cinema_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}

	if err := c.BindJSON(&ShowtimeReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := mh.ShowtimeService.AddShowtime(ShowtimeReq.BranchID, ShowtimeReq.MovieID, ShowtimeReq.TheaterID, ShowtimeReq.CinemaID, ShowtimeReq.StartTime, ShowtimeReq.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, "Add successfully")

}

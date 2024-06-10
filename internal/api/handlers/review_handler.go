package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	reviewService service.ReviewService
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService}
}

func (api *ReviewHandler) CreateReview(c *gin.Context) {
	var reviewRequest model.ReviewRequest

	if err := c.ShouldBind(&reviewRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		imageInfo, err := utils.UploadToCloudinary(file, header.Filename, "reviews")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		reviewRequest.ImageURL = imageInfo.ImageURL
		reviewRequest.ImageID = imageInfo.PublicID

	}

	if err := api.reviewService.CreateReview(reviewRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reviewRequest)
}

func (api *ReviewHandler) GetReviewsByMovieID(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	reviews, err := api.reviewService.GetReviewsByMovieID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

func (api *ReviewHandler) GetAverageRatingAndCountByMovieID(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("movie_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	avgRating, count, err := api.reviewService.GetAverageRatingAndCountByMovieID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"average_rating": avgRating,
		"review_count":   count,
	})
}

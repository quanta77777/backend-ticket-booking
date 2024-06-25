package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := uh.UserService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")
	user, err := uh.UserService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) CreateUser(c *gin.Context) {
	var userRequest model.UserRequest

	if err := c.ShouldBind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, header, err := c.Request.FormFile("image")
	if err == nil {
		defer file.Close()
		imageInfo, err := utils.UploadToCloudinary(file, header.Filename, "user-avatar")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		userRequest.ImageURL = imageInfo.ImageURL
		userRequest.ImageID = imageInfo.PublicID

	}

	err = uh.UserService.CreateUser(userRequest)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
	}

	c.JSON(http.StatusCreated, "Add successfully")
}

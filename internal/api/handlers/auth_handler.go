package handlers

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
	UserService *service.UserService
}

func (h *AuthHandler) Register(c *gin.Context) {
	var userRequest model.UserRequest

	if err := c.ShouldBind(&userRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := h.AuthService.HashPassword(userRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	userRequest.Password = hashedPassword

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

	if err := h.UserService.CreateUser(userRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	storedUser, err := h.UserService.GetUserByEmail(user.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !h.AuthService.CheckPasswordHash(user.Password, storedUser.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	tokenDetails, err := h.AuthService.CreateToken(storedUser.Email, storedUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var tokenRequest struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&tokenRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenDetails, err := h.AuthService.RefreshToken(tokenRequest.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenDetails)
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	email, _ := c.Get("email")
	role, _ := c.Get("role")

	c.JSON(http.StatusOK, gin.H{"email": email, "role": role})
}

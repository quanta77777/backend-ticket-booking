package handlers

import (
	"mime/multipart"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CinemaHandler struct {
	CinemaService *service.CinemaService
}

func NewCinemaHandler(cinemaService *service.CinemaService) *CinemaHandler {
	return &CinemaHandler{CinemaService: cinemaService}
}

func (cs *CinemaHandler) AddCinema(c *gin.Context) {

	filename, ok := c.Get("filePath")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename not found"})
	}

	file, ok := c.Get("file")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file not found"})
		return
	}
	imageInfo, err := utils.UploadToCloudinary(file.(multipart.File),
		filename.(string), "cinemalogo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var cinemaChainInput struct {
		Name  string                `form:"name" binding:"required"`
		Image *multipart.FileHeader `form:"file" binding:"required"`
	}
	if err := c.ShouldBind(&cinemaChainInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err = cs.CinemaService.AddCinema(cinemaChainInput.Name, imageInfo.ImageURL, imageInfo.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add cinema chain"})
	}

	newCinemaChain := model.Cinema{

		Name:     cinemaChainInput.Name,
		ImageURL: imageInfo.ImageURL,
		ImageID:  imageInfo.PublicID,
	}
	c.JSON(http.StatusCreated, newCinemaChain)
}
func (cs *CinemaHandler) UpdateCinema(c *gin.Context) {
	cinemaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}
	var cinemaRequest struct {
		Name  string                `form:"name"`
		Image *multipart.FileHeader `form:"file"`
	}
	if err := c.ShouldBind(&cinemaRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var imageURL string
	var publicID string

	if cinemaRequest.Image != nil {

		// Xác định xem bộ phim có hình ảnh trước đó không
		oldCinema, err := cs.CinemaService.GetCinemaByID(cinemaID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movie information"})
			return
		}
		oldPublicID := oldCinema.ImageID
		// Tiến hành tải lên hình ảnh mới lên Cloudinary
		file, err := cinemaRequest.Image.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer file.Close()

		imageInfo, err := utils.UploadToCloudinary(file, cinemaRequest.Image.Filename, "cinemalogo")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload image to Cloudinary"})
			return
		}

		err = utils.DeleteFromCloudinary(oldPublicID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image to Cloudinary"})
			return
		}
		imageURL = imageInfo.ImageURL
		publicID = imageInfo.PublicID
	}
	err = cs.CinemaService.UpdateCinema(cinemaID, cinemaRequest.Name, imageURL, publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	// Trả về kết quả thành công cho client
	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})

}
func (cs *CinemaHandler) GetAllCinema(c *gin.Context) {
	cinemas, err := cs.CinemaService.GetAllCinema()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinemas)

}

func (cs *CinemaHandler) GetCinemaByID(c *gin.Context) {
	cinemaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema ID"})
		return
	}
	cinema, err := cs.CinemaService.GetCinemaByID(cinemaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cinema)

}

func (cs *CinemaHandler) DeleteCinema(c *gin.Context) {
	cinemaID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema ID"})
		return
	}

	oldCinema, err := cs.CinemaService.GetCinemaByID(cinemaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get cinema  information"})
		return
	}
	oldPublicID := oldCinema.ImageID
	err = utils.DeleteFromCloudinary(oldPublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image to Cloudinary"})
		return
	}

	err = cs.CinemaService.DeleteCinema(cinemaID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cinema"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cinema delete successfully"})

}

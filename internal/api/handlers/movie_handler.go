package handlers

import (
	"mime/multipart"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MovieHandler struct {
	MovieService *service.MovieService
}

func NewMovieHandler(movieService *service.MovieService) *MovieHandler {
	return &MovieHandler{MovieService: movieService}
}

func (mh *MovieHandler) GetAllMovie(c *gin.Context) {
	status := c.Query("status")

	movies, err := mh.MovieService.GetAllMovie(status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (mh *MovieHandler) GetMovieByID(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}
	movies, err := mh.MovieService.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, movies)
}

func (mh *MovieHandler) AddMovie(c *gin.Context) {

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
		filename.(string), "movies")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var movieRequest struct {
		Title       string                `form:"title" binding:"required"`
		Director    string                `form:"director" binding:"required"`
		Genre       string                `form:"genre" binding:"required"`
		Duration    int                   `form:"duration" binding:"required"`
		Image       *multipart.FileHeader `form:"file" binding:"required"`
		Description string                `form:"description" binding:"required"`
		ReleaseDate string                `form:"release_date" binding:"required"`
		EndDate     string                `form:"end_date" binding:"required"`
	}
	if err := c.ShouldBind(&movieRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie, err := mh.MovieService.AddMovie(movieRequest.Title, movieRequest.Director, movieRequest.Genre, movieRequest.Duration, imageInfo.ImageURL, imageInfo.PublicID, movieRequest.Description, movieRequest.ReleaseDate, movieRequest.EndDate)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusCreated, movie)

}

func (mh *MovieHandler) UpdateMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	var movieRequest struct {
		Title    string                `form:"title"`
		Director string                `form:"director"`
		Genre    string                `form:"genre" `
		Duration int                   `form:"duration" `
		Image    *multipart.FileHeader `form:"file"`
	}

	if err := c.ShouldBind(&movieRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	var imageURL string
	var publicID string

	if movieRequest.Image != nil {

		// Xác định xem bộ phim có hình ảnh trước đó không
		oldMovie, err := mh.MovieService.GetMovieByID(movieID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movie information"})
			return
		}
		oldPublicID := oldMovie.ImageID
		// Tiến hành tải lên hình ảnh mới lên Cloudinary
		file, err := movieRequest.Image.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer file.Close()

		imageInfo, err := utils.UploadToCloudinary(file, movieRequest.Image.Filename, "movies")
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

	err = mh.MovieService.UpdateMovie(movieID, movieRequest.Title, movieRequest.Director, movieRequest.Genre, movieRequest.Duration, imageURL, publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	// Trả về kết quả thành công cho client
	c.JSON(http.StatusOK, gin.H{"message": "Movie updated successfully"})

}

func (mh *MovieHandler) DeleteMovie(c *gin.Context) {
	movieID, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movie ID"})
		return
	}

	oldMovie, err := mh.MovieService.GetMovieByID(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get movie information"})
		return
	}
	oldPublicID := oldMovie.ImageID
	err = utils.DeleteFromCloudinary(oldPublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image to Cloudinary"})
		return
	}

	err = mh.MovieService.DeleteMovie(movieID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie delete successfully"})
}

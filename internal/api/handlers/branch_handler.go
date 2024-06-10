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

type BranchHandler struct {
	BranchService *service.BranchService
}

func NewBranchHandler(branchService *service.BranchService) *BranchHandler {
	return &BranchHandler{BranchService: branchService}
}

func (bh *BranchHandler) AddBranch(c *gin.Context) {

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
		filename.(string), "branchlogo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var branchReq struct {
		CinemaId int                   `form:"cinema_id" binding:"required"`
		Name     string                `form:"name" binding:"required"`
		Address  string                `form:"address" binding:"required"`
		Image    *multipart.FileHeader `form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&branchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	BranchID, err := bh.BranchService.AddBranch(branchReq.CinemaId, branchReq.Name, branchReq.Address, imageInfo.ImageURL, imageInfo.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add cinema chain"})
	}

	newBranch := model.Branch{
		BranchId: int(BranchID),
		Name:     branchReq.Name,
		Address:  branchReq.Address,
		CinemaId: branchReq.CinemaId,
		ImageURL: imageInfo.ImageURL,
		ImageID:  imageInfo.PublicID,
	}
	c.JSON(http.StatusCreated, newBranch)

}

func (bh *BranchHandler) UpdateCinema(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}

	var branchReq struct {
		Name    string                `form:"name"`
		Address string                `form:"address"`
		Image   *multipart.FileHeader `form:"file"`
	}
	if err := c.ShouldBind(&branchReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var imageURL string
	var publicID string

	if branchReq.Image != nil {

		// Xác định xem bộ phim có hình ảnh trước đó không
		oldBranch, err := bh.BranchService.GetBranchByID(branchID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch information"})
			return
		}
		oldPublicID := oldBranch.ImageID
		// Tiến hành tải lên hình ảnh mới lên Cloudinary
		file, err := branchReq.Image.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open image file"})
			return
		}
		defer file.Close()

		imageInfo, err := utils.UploadToCloudinary(file, branchReq.Image.Filename, "branchlogo")
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
	err = bh.BranchService.UpdateBranch(branchID, branchReq.Name, branchReq.Address, imageURL, publicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update branch"})
		return
	}

	// Trả về kết quả thành công cho client
	c.JSON(http.StatusOK, gin.H{"message": "Branch updated successfully"})

}

func (bh *BranchHandler) GetAllBranch(c *gin.Context) {
	branches, err := bh.BranchService.GetAllBranch()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branches)
}

func (bh *BranchHandler) GetBranchByID(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}
	branche, err := bh.BranchService.GetBranchByID(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branche)
}

func (bh *BranchHandler) GetBranchesByCinemaChainID(c *gin.Context) {
	cinemaChainIDStr := c.Param("id")
	cinemaChainID, err := strconv.Atoi(cinemaChainIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cinema chain ID"})
		return
	}

	branches, err := bh.BranchService.GetBranchesByCinemaChainID(cinemaChainID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, branches)
}

func (bh *BranchHandler) DeleteBranch(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid branch ID"})
		return
	}
	oldBranch, err := bh.BranchService.GetBranchByID(branchID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get branch  information"})
		return
	}
	oldPublicID := oldBranch.ImageID
	err = utils.DeleteFromCloudinary(oldPublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete image to Cloudinary"})
		return
	}

	err = bh.BranchService.DeleteBranch(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete branch"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Branch delete successfully"})

}

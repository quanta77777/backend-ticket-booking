package handlers

import (
	"mime/multipart"
	"movie-ticket-booking/internal/service"
	"movie-ticket-booking/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	ProductHandler *service.ProductService
}

func NewProductHandler(productHandle *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductHandler: productHandle}
}

func (ph *ProductHandler) CreateProductForBranch(c *gin.Context) {
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
		filename.(string), "products")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var productReq struct {
		BranchId    int                   `form:"branch_id" binding:"required"`
		Name        string                `form:"name" binding:"required"`
		Description string                `form:"description" binding:"required"`
		Price       int64                 `form:"price" binding:"required"`
		Image       *multipart.FileHeader `form:"file" binding:"required"`
	}

	if err := c.ShouldBind(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ph.ProductHandler.CreateProductForBranch(productReq.BranchId, productReq.Name, productReq.Description, productReq.Price, imageInfo.ImageURL, imageInfo.PublicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Add successful"})
}

func (ph *ProductHandler) GetProductByBranchID(c *gin.Context) {
	branchID, err := strconv.Atoi(c.Query("branch-id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Branch id"})
		return
	}

	products, err := ph.ProductHandler.GetProductByBranchID(branchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

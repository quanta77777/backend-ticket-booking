package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type ProductService struct {
	ProductService *repository.ProductRepository
}

func NewProductService(db *sql.DB) *ProductService {
	return &ProductService{ProductService: repository.NewProductRepository(db)}
}

func (ps *ProductService) CreateProductForBranch(branchID int, name, description string, price int64, ImageURL, ImageID string) error {
	return ps.ProductService.CreateProductForBranch(branchID, name, description, price, ImageURL, ImageID)
}

func (ps *ProductService) GetProductByBranchID(branchId int) ([]model.Product, error) {
	return ps.ProductService.GetProductByBranchID(branchId)
}

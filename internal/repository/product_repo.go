package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db}
}

func (pr *ProductRepository) CreateProductForBranch(branchID int, name, description string, price int64, ImageURL, ImageID string) error {
	query := `INSERT INTO product(branch_id, name, description, price, image_url, image_id) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := pr.DB.Exec(query, branchID, name, description, price, ImageURL, ImageID)
	return err
}

func (pr *ProductRepository) GetProductByBranchID(branchID int) ([]model.Product, error) {
	rows, err := pr.DB.Query("SELECT * FROM product WHERE branch_id = ?", branchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var products []model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(&product.ProductID, &product.BranchId, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.ImageID)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

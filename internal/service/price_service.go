package service

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type PriceService struct {
	PriceService *repository.PriceRepository
}

func NewPriceService(db *sql.DB) *PriceService {
	return &PriceService{PriceService: repository.NewPriceRepository(db)}
}

func (ps *PriceService) CreatePriceForShowtime(req model.Price) error {
	return ps.PriceService.CreatePriceForShowtime(req)
}

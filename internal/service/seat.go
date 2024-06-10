package service

import (
	"database/sql"
	"movie-ticket-booking/internal/repository"

	"movie-ticket-booking/internal/model"
)

type SeatService struct {
	SeatService *repository.SeatRopository
	priceRepo   *repository.PriceRepository
}

func NewSeatService(db *sql.DB) *SeatService {
	return &SeatService{SeatService: repository.NewSeatRopository(db), priceRepo: repository.NewPriceRepository(db)}
}

func (ss *SeatService) AddSeat(theaterId int, seatType, seatNumber string) error {
	return ss.SeatService.AddSeat(theaterId, seatType, seatNumber)
}

func (ss *SeatService) GetSeatByTheaterID(theaterId int) ([]model.Seat, error) {
	return ss.SeatService.GetSeatByTheaterID(theaterId)
}

func (ss *SeatService) GetSeatsWithReservationStatusAndPrices(theaterID, showtimeID int) ([]model.SeatWithReservation, map[string]int64, error) {
	// return ss.SeatService.GetSeatsWithReservationStatus(theaterID, showtimeID)
	seats, err := ss.SeatService.GetSeatsWithReservationStatus(theaterID, showtimeID)
	if err != nil {
		return nil, nil, err
	}
	prices, err := ss.priceRepo.GetPricesByShowtimeID(showtimeID)
	if err != nil {
		return nil, nil, err
	}

	return seats, prices, nil
}

func (ss *SeatService) CreateReservation(reservation model.SeatReservation) error {
	return ss.SeatService.CreateReservation(reservation)
}

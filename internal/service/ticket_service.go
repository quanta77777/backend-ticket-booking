package service

import (
	"movie-ticket-booking/internal/model"
	"movie-ticket-booking/internal/repository"
)

type TicketService struct {
	ticketRepo *repository.TicketRepository
}

func NewTicketService(ticketRepo *repository.TicketRepository) *TicketService {
	return &TicketService{ticketRepo: ticketRepo}
}

func (s *TicketService) CreateTicket(ticket *model.Ticket) (int, error) {
	return s.ticketRepo.CreateTicket(ticket)
}
func (s *TicketService) AddProductWithTicketId(ticket *model.TicketProduct) error {
	return s.ticketRepo.AddProductWithTicketId(ticket)
}
func (s *TicketService) AddSeatWithTicketId(ticket *model.TicketSeat) error {
	return s.ticketRepo.AddSeatWithTicketId(ticket)
}
func (s *TicketService) UserHasTicketForMovie(userID, movieID int) (bool, error) {
	return s.ticketRepo.UserHasTicketForMovie(userID, movieID)
}

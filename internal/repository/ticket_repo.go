package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
)

type TicketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) *TicketRepository {
	return &TicketRepository{db: db}
}

func (r *TicketRepository) CreateTicket(ticket *model.Ticket) (int, error) {
	query := `
	INSERT INTO ticket (user_id, showtime_id, movie_id, price)
	VALUES (?, ?, ?, ?)
`

	result, err := r.db.Exec(query, ticket.UserID, ticket.ShowtimeID, ticket.MovieID, ticket.Price)
	if err != nil {
		return 0, err
	}

	ticketID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(ticketID), nil
}

func (r *TicketRepository) AddProductWithTicketId(ticket *model.TicketProduct) error {
	query := `
	INSERT INTO ticket_product (ticket_id, product_id)
	VALUES (?, ?)
`

	_, err := r.db.Exec(query, ticket.TicketID, ticket.ProductID)
	if err != nil {
		return err
	}

	return nil
}
func (r *TicketRepository) AddSeatWithTicketId(ticket *model.TicketSeat) error {
	query := `
	INSERT INTO ticket_seat (ticket_id, seat_id)
	VALUES (?, ?)
`

	_, err := r.db.Exec(query, ticket.TicketID, ticket.SeatID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TicketRepository) UserHasTicketForMovie(userID, movieID int) (bool, error) {
	query := `
        SELECT COUNT(*) 
        FROM ticket
        WHERE user_id = ? AND movie_id = ?
    `
	var count int
	err := r.db.QueryRow(query, userID, movieID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

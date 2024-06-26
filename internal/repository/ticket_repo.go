package repository

import (
	"database/sql"
	"movie-ticket-booking/internal/model"
	"time"
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

func (r *TicketRepository) GetTicketByUserId(userID int) ([]model.Ticket, error) {
	query := `
        SELECT *
        FROM ticket
        WHERE user_id = ? 
    `
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket
	for rows.Next() {
		var ticket model.Ticket
		var createdAtBytes []byte
		if err := rows.Scan(&ticket.TicketID, &ticket.UserID, &ticket.ShowtimeID, &ticket.MovieID, &ticket.Price, &createdAtBytes); err != nil {
			return nil, err
		}
		createdAt, err := time.Parse("2006-01-02 15:04:05", string(createdAtBytes))
		if err != nil {
			return nil, err
		}
		ticket.CreatedAt = createdAt
		tickets = append(tickets, ticket)
	}
	return tickets, nil

}

func (r *TicketRepository) GetSeatByTicketId(ticketID int) ([]model.TicketSeat, error) {
	query := `
        SELECT *
        FROM ticket_seat
        WHERE ticket_id = ? 
    `
	rows, err := r.db.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.TicketSeat
	for rows.Next() {
		var ticket model.TicketSeat
		if err := rows.Scan(&ticket.TicketSeatID, &ticket.TicketID, &ticket.SeatID); err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}
	return tickets, nil

}

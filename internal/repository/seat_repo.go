package repository

import (
	"database/sql"
	"fmt"

	"errors"
	"movie-ticket-booking/internal/model"
)

type SeatRopository struct {
	DB *sql.DB
}

func NewSeatRopository(db *sql.DB) *SeatRopository {
	return &SeatRopository{DB: db}
}

func (sr *SeatRopository) AddSeat(theaterId int, seatType, seatNumber string) error {
	query := "INSERT INTO seat (theater_id, seat_number, seat_type) VALUES (?, ?, ?)"
	_, err := sr.DB.Exec(query, theaterId, seatNumber, seatType)

	if err != nil {
		return errors.New("failed to insert seat")
	}
	return nil

}

func (sr *SeatRopository) GetSeatByTheaterID(theaterID int) ([]model.Seat, error) {
	rows, err := sr.DB.Query("SELECT * FROM seat WHERE theater_id = ?", theaterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var seats []model.Seat
	for rows.Next() {
		var seat model.Seat
		err := rows.Scan(&seat.SeatId, &seat.TheaterId, &seat.SeatNumber, &seat.SeatType)
		if err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}
func (sr *SeatRopository) GetSeatsWithReservationStatus(theaterID, showtimeID int) ([]model.SeatWithReservation, error) {
	var seats []model.SeatWithReservation
	query := `SELECT s.seat_id, s.theater_id, s.seat_number, s.seat_type,
                     IF(sr.seat_id IS NOT NULL, TRUE, FALSE) AS is_reserved
              FROM seat s
              LEFT JOIN seat_reservation sr ON s.seat_id = sr.seat_id AND sr.showtime_id = ?
              WHERE s.theater_id = ?`
	rows, err := sr.DB.Query(query, showtimeID, theaterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var seat model.SeatWithReservation
		if err := rows.Scan(&seat.SeatId, &seat.TheaterId, &seat.SeatNumber, &seat.SeatType, &seat.IsReserved); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return seats, nil
}

func (sr *SeatRopository) CreateReservation(reservation model.SeatReservation) error {
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM seat_reservation WHERE seat_id = ? AND showtime_id = ?)`
	err := sr.DB.QueryRow(checkQuery, reservation.SeatID, reservation.ShowtimeID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if reservation exists: %w", err)
	}
	if exists {
		return errors.New("reservation already exists")
	}

	// Tạo mới bản ghi nếu không tồn tại
	insertQuery := `INSERT INTO seat_reservation (seat_id, showtime_id) VALUES (?, ?)`
	_, err = sr.DB.Exec(insertQuery, reservation.SeatID, reservation.ShowtimeID)
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}
	return nil
}

func (sr *SeatRopository) GetSeatByTicketId(ticketID int) ([]model.TicketSeat, error) {
	query := `
        SELECT *
        FROM ticket_seat
        WHERE ticket_id = ? 
    `
	rows, err := sr.DB.Query(query, ticketID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []model.TicketSeat
	for rows.Next() {
		var seat model.TicketSeat
		if err := rows.Scan(&seat.TicketSeatID, &seat.TicketID, &seat.SeatID); err != nil {
			return nil, err
		}

		seats = append(seats, seat)
	}
	return seats, nil

}

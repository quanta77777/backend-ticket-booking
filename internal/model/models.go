package model

import "time"

type Branch struct {
	BranchId  int       `json:"branch_id"`
	CinemaId  int       `json:"cinema_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	ImageURL  string    `json:"image_url"`
	ImageID   string    `json:"image_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Cinema struct {
	CinemaID  int       `json:"cinema_id"`
	Name      string    `json:"name"`
	ImageURL  string    `json:"image_url"`
	ImageID   string    `json:"image_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Showtime struct {
	ShowtimeID int       `json:"showtime_id"`
	CinemaID   int       `json:"cinema_id"`
	BranchID   int       `json:"branch_id"`
	TheaterID  int       `json:"theater_id"`
	MovieID    int       `json:"movie_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	CreatedAt  time.Time `json:"created_at"`
}

type Movie struct {
	MovieID     int       `json:"movie_id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	Genre       string    `json:"genre"`
	Duration    int       `json:"duration"`
	ImageURL    string    `json:"image_url"`
	ImageID     string    `json:"public_id"`
	Discription string    `json:"discription"`
	ReleaseDate string    `json:"release_date"`
	EndDate     string    `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
}
type MovieRequest struct {
	Title     string    `json:"title"`
	Director  string    `json:"director"`
	Genre     string    `json:"genre"`
	Duration  int       `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

type Price struct {
	ShowtimeID int       `json:"showtime_id"`
	SeatType   string    `json:"seat_type"`
	Price      int64     `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

type Product struct {
	ProductID   int       `json:"product_id"`
	BranchId    int       `json:"branch_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	ImageURL    string    `json:"image_url"`
	ImageID     string    `json:"public_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type Review struct {
	ReviewID  int    `json:"review_id"`
	UserID    int    `json:"user_id"`
	MovieID   int    `json:"movie_id"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
	ImageURL  string `json:"image_url"`
	ImageID   string `json:"image_id"`
	CreatedAt time.Time
}
type ReviewRequest struct {
	UserID   int    `form:"user_id" binding:"required"`
	MovieID  int    `form:"movie_id" binding:"required"`
	Rating   int    `form:"rating" binding:"required"`
	Comment  string `form:"comment" binding:"required"`
	ImageURL string `form:"image_url" `
	ImageID  string `form:"image_id"`
}
type Seat struct {
	SeatId     int       `json:"seat_id"`
	TheaterId  int       `json:"theater_id"`
	SeatNumber string    `json:"seat_number"`
	SeatType   string    `json:"seat_type"`
	CreatedAt  time.Time `json:"created_at"`
}

type SeatReservation struct {
	SeatID     int `json:"seat_id"`
	ShowtimeID int `json:"showtime_id"`
	// CreatedAt  time.Time `json:"created_at"`
}

type SeatWithReservation struct {
	Seat
	IsReserved bool `json:"is_reserved"`
}

type Theater struct {
	TheaterId int       `json:"theater_id"`
	BranchId  int       `json:"branch_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type Ticket struct {
	TicketID   int       `json:"ticket_id"`
	UserID     int       `json:"user_id"`
	MovieID    int       `json:"movie_id"`
	ShowtimeID int       `json:"showtime_id"`
	Price      int       `json:"price"`
	CreatedAt  time.Time `json:"created_at"`
}

type TicketProduct struct {
	TicketProductID int `json:"ticket_product_id"`
	TicketID        int `json:"ticket_id"`
	ProductID       int `json:"product_id"`
}

type TicketSeat struct {
	TicketSeatID int `json:"ticket_seat_id"`
	TicketID     int `json:"ticket_id"`
	SeatID       int `json:"seat_id"`
}

type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

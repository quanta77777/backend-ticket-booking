package routes

import (
	"database/sql"
	"net/http"

	"movie-ticket-booking/internal/api/handlers"
	"movie-ticket-booking/internal/middleware"
	"movie-ticket-booking/internal/repository"
	"movie-ticket-booking/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(db *sql.DB) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize repositories, services and handlers
	userRepo := &repository.UserRepository{DB: db}
	userService := service.NewUserService(db)
	userHandler := handlers.NewUserHandler(userService)

	authService := &service.AuthService{UserRepo: userRepo}
	authHandler := &handlers.AuthHandler{
		AuthService: authService,
		UserService: userService,
	}

	theaterService := service.NewTheaterService(db)
	theaterHandler := handlers.NewTheaterHandler(theaterService)

	cinemaService := service.NewCinemaService(db)
	cinemaHandler := handlers.NewCinemaHandler(cinemaService)

	branchService := service.NewBranchService(db)
	branchHandler := handlers.NewBranchHandler(branchService)

	showtimeService := service.NewShowtimeService(db)
	showtimeHandler := handlers.NewShowtimeHandler(showtimeService)

	movieService := service.NewMovieService(db)
	movieHandler := handlers.NewMovieHandler(movieService)

	seatService := service.NewSeatService(db)
	seatHandler := handlers.NewSeatHandler(seatService)

	priceService := service.NewPriceService(db)
	priceHandler := handlers.NewPriceHandler(priceService)

	productService := service.NewProductService(db)
	productHandler := handlers.NewProductHandler(productService)

	reviewRepo := repository.NewReviewRepository(db)
	reviewService := service.NewReviewService(reviewRepo)
	reviewHandler := handlers.NewReviewHandler(reviewService)

	ticketRepo := repository.NewTicketRepository(db)
	ticketService := service.NewTicketService(ticketRepo)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	// Define routes
	router.GET("/api/users", userHandler.GetAllUsers)
	router.GET("/api/users/:id", userHandler.GetUserByID)
	router.POST("/api/users", userHandler.CreateUser)

	router.POST("/refresh", authHandler.Refresh)
	router.POST("/api/auth/register", authHandler.Register)
	router.POST("/api/auth/login", authHandler.Login)

	router.POST("/api/cinema", middleware.FileUploadMiddleware(), cinemaHandler.AddCinema)
	router.GET("/api/cinema", cinemaHandler.GetAllCinema)
	router.GET("/api/cinema/:id", cinemaHandler.GetCinemaByID)
	router.PUT("/api/cinema/:id", cinemaHandler.UpdateCinema)
	router.DELETE("/api/cinema/:id", cinemaHandler.DeleteCinema)

	router.POST("/api/branch", middleware.FileUploadMiddleware(), branchHandler.AddBranch)
	router.GET("/api/branch/cinema/:id", branchHandler.GetBranchesByCinemaChainID)
	router.GET("/api/branch", branchHandler.GetAllBranch)
	router.GET("/api/branch/:id", branchHandler.GetBranchByID)
	router.PUT("/api/branch/:id", branchHandler.UpdateCinema)
	router.DELETE("/api/branch/:id", branchHandler.DeleteBranch)

	router.GET("/api/theater/branch/:id", theaterHandler.GetAllTheaterByBranchID)
	router.GET("/api/theater/:id", theaterHandler.GetTheaterByID)
	router.POST("/api/theater", theaterHandler.AddTheater)

	router.POST("/api/seat", seatHandler.AddSeat)
	router.GET("/api/seat/theater/:id", seatHandler.GetSeatBytheaterID)
	router.GET("/api/seats", seatHandler.GetSeatsWithReservationStatusAndPrices)
	router.POST("/api/reserve-seat", seatHandler.ReserveSeat)

	router.POST("/api/price", priceHandler.CreatePriceForShowtime)

	router.POST("/api/product", middleware.FileUploadMiddleware(), productHandler.CreateProductForBranch)
	router.GET("/api/product", productHandler.GetProductByBranchID)

	router.POST("/api/showtime", showtimeHandler.AddShowtime)
	router.GET("/api/showtime", showtimeHandler.GetShowtimeByDay)
	router.GET("/api/showtime/branch/:id", showtimeHandler.GetShowtimeWithBranch)
	router.GET("/api/showtime/cinema/:id", showtimeHandler.GetShowtimeWithCinema)

	router.GET("/api/movie", movieHandler.GetAllMovie)
	router.GET("/api/movie/:id", movieHandler.GetMovieByID)
	router.POST("/api/movie", middleware.FileUploadMiddleware(), movieHandler.AddMovie)
	router.PUT("/api/movie/:id", movieHandler.UpdateMovie)
	router.DELETE("/api/movie/:id", movieHandler.DeleteMovie)

	router.POST("/api/ticket", ticketHandler.CreateTicket)
	router.POST("/api/ticket/product", ticketHandler.AddProductWithTicketId)
	router.POST("/api/ticket/seat", ticketHandler.AddSeatWithTicketId)
	router.GET("/api/ticket/user/:user_id/movie/:movie_id", ticketHandler.UserHasTicketForMovie)

	router.POST("/api/movies/review", reviewHandler.CreateReview)
	router.GET("/api/movies/:movie_id/reviews", reviewHandler.GetReviewsByMovieID)
	router.GET("/api/movies/:movie_id/average_rating", reviewHandler.GetAverageRatingAndCountByMovieID)

	authorized := router.Group("/api/admin")
	authorized.Use(middleware.JWTAuthMiddleware(authService), middleware.CheckRole("admin"))
	authorized.GET("/dashboard", func(c *gin.Context) {

		c.JSON(http.StatusOK, gin.H{"message": "Welcome to admin dashboard!"})
	})

	profile := router.Group("/profile")
	profile.Use(middleware.JWTAuthMiddleware(authService))
	profile.GET("", authHandler.GetProfile)

	return router
}

package http

import (
	"net/http"

	internal_db "github.com/nambuitechx/nam-tcp/internal/db"
	"github.com/nambuitechx/nam-tcp/internal/handlers"
	"github.com/nambuitechx/nam-tcp/internal/repositories"
	"github.com/nambuitechx/nam-tcp/internal/services"
)

type HttpServer struct {
	Mux *http.ServeMux
}

func NewHttpServer() *HttpServer {
	db := internal_db.GetDBConnection()
	internal_db.Up(db)

	userRepository := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepository)
	userHandler := handlers.NewUserHandler(userService)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.WriteHealth)
	mux.HandleFunc("/", handlers.WriteHealth)

	mux.HandleFunc("GET /api/v1/users", userHandler.GetUsers())
	mux.HandleFunc("POST /api/v1/users", userHandler.CreateUser())

	return &HttpServer{
		Mux: mux,
	}
}

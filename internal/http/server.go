package http

import (
	"net/http"

	"github.com/nambuitechx/nam-tcp/internal/handlers"
	"github.com/nambuitechx/nam-tcp/internal/deps"
)

type HttpServer struct {
	Mux *http.ServeMux
}

func NewHttpServer(deps *deps.ServiceDeps) *HttpServer {
	userHandler := handlers.NewUserHandler(deps.UserService)
	targetHandler := handlers.NewTargetHandler(deps.TargetService)
	userPATHandler := handlers.NewUserPATHandler(deps.UserPATService)

	mux := http.NewServeMux()

	mux.HandleFunc("/health", handlers.WriteHealth)
	mux.HandleFunc("/", handlers.WriteHealth)

	mux.HandleFunc("GET /api/v1/users", userHandler.GetUsers())
	mux.HandleFunc("POST /api/v1/users", userHandler.CreateUser())

	mux.HandleFunc("GET /api/v1/targets", targetHandler.GetTargets())
	mux.HandleFunc("POST /api/v1/targets", targetHandler.CreateTarget())

	mux.HandleFunc("GET /api/v1/user-pats", userPATHandler.GetUserPATs())
	mux.HandleFunc("POST /api/v1/user-pats", userPATHandler.CreateUserPAT())
	mux.HandleFunc("DELETE /api/v1/user-pats/{id}", userPATHandler.RevokeUserPAT())
	mux.HandleFunc("DELETE /api/v1/user-pats/expired", userPATHandler.RevokeExpiredUserPATs())

	return &HttpServer{
		Mux: mux,
	}
}

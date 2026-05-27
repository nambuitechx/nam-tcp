package deps

import (
	"database/sql"

	"github.com/nambuitechx/nam-tcp/internal/repositories"
	"github.com/nambuitechx/nam-tcp/internal/services"
)

type ServiceDeps struct {
	UserRepository *repositories.UserRepository
	TargetRepository *repositories.TargetRepository
	UserPATRepository *repositories.UserPATRepository
	UserService *services.UserService
	TargetService *services.TargetService
	UserPATService *services.UserPATService
}

func NewServiceDeps(db *sql.DB) *ServiceDeps {
	userRepository := repositories.NewUserRepository(db)
	targetRepository := repositories.NewTargetRepository(db)
	userPATRepository := repositories.NewUserPATRepository(db)

	userService := services.NewUserService(userRepository)
	targetService := services.NewTargetService(targetRepository)
	userPATService := services.NewUserPATService(userPATRepository)

	return &ServiceDeps{
		UserRepository: userRepository,
		TargetRepository: targetRepository,
		UserPATRepository: userPATRepository,
		UserService: userService,
		TargetService: targetService,
		UserPATService: userPATService,
	}
}
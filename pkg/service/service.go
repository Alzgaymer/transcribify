package service

import (
	"log"
	"os"
	"transcribify/pkg/auth"
	"transcribify/pkg/finders"
	"transcribify/pkg/hash"
	"transcribify/pkg/repository"
)

type (
	Services struct {
		Manager       auth.TokenManager
		Authorization auth.Authorization
		Finder        finders.Finder
	}
)

func New(repository repository.Repository, finder finders.Finder, hasher hash.PasswordHasher) *Services {
	manager, err := auth.NewManager(os.Getenv("JWT_SALT"))
	if err != nil {
		log.Fatal(err)
	}

	return &Services{
		Manager:       manager,
		Authorization: auth.NewAuthorizationManager(repository.User, manager, hasher),
		Finder:        finder,
	}
}

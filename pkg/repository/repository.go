package repository

import (
	"context"
	"github.com/jackc/pgx/v5"
	"transcribify/internal/models"
)

type (
	Repository struct {
		Video Video
		User  User
	}

	Video interface {
		Create(context.Context, models.YTVideo, models.VideoRequest) (int, error)
		Read(context.Context, models.VideoRequest) (models.YTVideo, error)
		Update(context.Context, models.VideoRequest, models.YTVideo) error
		Delete(context.Context, models.VideoRequest) error
	}

	User interface {
	}
)

func NewRepositories(client *pgx.Conn) *Repository {
	return &Repository{
		Video: NewYTVideoRepository(client),
		User:  nil,
	}
}

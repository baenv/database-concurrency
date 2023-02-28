package gen

import (
	"context"
	"database-concurrency/internal/repository"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	CreateTicketID(ctx context.Context, uniqueID uuid.UUID) (string, error)
}

func New(repo repository.Repository, log *logrus.Logger) Controller {
	return gen{
		repo: repo,
		log:  log,
	}
}

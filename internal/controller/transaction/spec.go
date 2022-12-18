package transaction

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/repository"

	"github.com/sirupsen/logrus"
)

type Controller interface {
	OneByHash(ctx context.Context, hash string) (*ent.Transaction, error)
}

func New(repo repository.Repositoy, log *logrus.Logger) Controller {
	return transaction{
		repo: repo,
		log:  log,
	}
}

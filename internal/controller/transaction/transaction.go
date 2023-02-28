package transaction

import (
	"context"
	"database-concurrency/ent"
	"database-concurrency/internal/repository"

	"github.com/sirupsen/logrus"
)

type transaction struct {
	repo repository.Repository
	log  *logrus.Logger
}

func (t transaction) OneByHash(ctx context.Context, hash string) (*ent.Transaction, error) {
	return t.repo.Transaction().OneByHash(ctx, hash)
}

package user

import (
	"database-concurrency/ent"
)

type pg struct {
	client *ent.UserClient
}

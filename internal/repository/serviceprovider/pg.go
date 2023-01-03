package serviceprovider

import (
	"database-concurrency/ent"
)

type pg struct {
	client *ent.ServiceProdiverClient
}

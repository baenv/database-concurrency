package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// CreateTicketLog holds the schema definition for the CreateTicketLog entity.
type CreateTicketLog struct {
	ent.Schema
}

// Fields of the CreateTicketLog.
func (CreateTicketLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New).StorageKey("ticket_id"),
		field.UUID("unique_id", uuid.UUID{}).Unique(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the CreateTicketLog.
func (CreateTicketLog) Edges() []ent.Edge {
	return []ent.Edge{}
}

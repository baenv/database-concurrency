package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TicketEvent holds the schema definition for the TicketEvent entity.
type TicketEvent struct {
	ent.Schema
}

// Fields of the TicketEvent.
func (TicketEvent) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("ticket_id", uuid.UUID{}),
		field.UUID("user_id", uuid.UUID{}),
		field.String("type"),
		field.JSON("metadada", map[string]interface{}{}),
		field.String("versions"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the TicketEvent.
func (TicketEvent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("ticket_events").
			Field("user_id").
			Required().
			Unique(),
		edge.From("ticket", Ticket.Type).
			Ref("ticket_events").
			Field("ticket_id").
			Required().
			Unique(),
	}
}

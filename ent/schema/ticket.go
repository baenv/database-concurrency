package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Ticket holds the schema definition for the Ticket entity.
type Ticket struct {
	ent.Schema
}

// Fields of the Ticket.
func (Ticket) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("status"),
		field.UUID("user_id", uuid.UUID{}),
		field.JSON("metadata", map[string]interface{}{}).Optional(),
		field.String("versions").Optional(),
		field.UUID("last_event_id", uuid.UUID{}).Optional(),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the Ticket.
func (Ticket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("tickets").
			Field("user_id").
			Required().
			Unique(),
		edge.To("last_event", TicketEvent.Type).
			Field("last_event_id").
			Unique(),
		edge.To("ticket_events", TicketEvent.Type).
			Annotations(
				entsql.Annotation{
					OnDelete: entsql.Cascade,
				},
			),
	}
}

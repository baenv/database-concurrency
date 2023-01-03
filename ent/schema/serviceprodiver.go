package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ServiceProdiver holds the schema definition for the ServiceProdiver entity.
type ServiceProdiver struct {
	ent.Schema
}

// Fields of the ServiceProdiver.
func (ServiceProdiver) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.String("name"),
		field.String("email"),
		field.String("phone"),
		field.String("verdor_ref"),
		field.Time("created_at").Default(time.Now),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the ServiceProdiver.
func (ServiceProdiver) Edges() []ent.Edge {
	return nil
}

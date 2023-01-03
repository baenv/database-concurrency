// Code generated by ent, DO NOT EDIT.

package ticketevent

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the ticketevent type in the database.
	Label = "ticket_event"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTicketID holds the string denoting the ticket_id field in the database.
	FieldTicketID = "ticket_id"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// FieldMetadada holds the string denoting the metadada field in the database.
	FieldMetadada = "metadada"
	// FieldVersions holds the string denoting the versions field in the database.
	FieldVersions = "versions"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeTicket holds the string denoting the ticket edge name in mutations.
	EdgeTicket = "ticket"
	// Table holds the table name of the ticketevent in the database.
	Table = "ticket_events"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "ticket_events"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// TicketTable is the table that holds the ticket relation/edge.
	TicketTable = "ticket_events"
	// TicketInverseTable is the table name for the Ticket entity.
	// It exists in this package in order to avoid circular dependency with the "ticket" package.
	TicketInverseTable = "tickets"
	// TicketColumn is the table column denoting the ticket relation/edge.
	TicketColumn = "ticket_id"
)

// Columns holds all SQL columns for ticketevent fields.
var Columns = []string{
	FieldID,
	FieldTicketID,
	FieldUserID,
	FieldType,
	FieldMetadada,
	FieldVersions,
	FieldCreatedAt,
	FieldUpdatedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

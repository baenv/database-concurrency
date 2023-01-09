// Code generated by ent, DO NOT EDIT.

package ticket

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the ticket type in the database.
	Label = "ticket"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldStatus holds the string denoting the status field in the database.
	FieldStatus = "status"
	// FieldUserID holds the string denoting the user_id field in the database.
	FieldUserID = "user_id"
	// FieldMetadata holds the string denoting the metadata field in the database.
	FieldMetadata = "metadata"
	// FieldVersions holds the string denoting the versions field in the database.
	FieldVersions = "versions"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeTicketEvents holds the string denoting the ticket_events edge name in mutations.
	EdgeTicketEvents = "ticket_events"
	// Table holds the table name of the ticket in the database.
	Table = "tickets"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "tickets"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_id"
	// TicketEventsTable is the table that holds the ticket_events relation/edge.
	TicketEventsTable = "ticket_events"
	// TicketEventsInverseTable is the table name for the TicketEvent entity.
	// It exists in this package in order to avoid circular dependency with the "ticketevent" package.
	TicketEventsInverseTable = "ticket_events"
	// TicketEventsColumn is the table column denoting the ticket_events relation/edge.
	TicketEventsColumn = "ticket_id"
)

// Columns holds all SQL columns for ticket fields.
var Columns = []string{
	FieldID,
	FieldStatus,
	FieldUserID,
	FieldMetadata,
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
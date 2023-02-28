// Code generated by ent, DO NOT EDIT.

package createticketlog

import (
	"time"

	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the createticketlog type in the database.
	Label = "create_ticket_log"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "ticket_id"
	// FieldUniqueID holds the string denoting the unique_id field in the database.
	FieldUniqueID = "unique_id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// Table holds the table name of the createticketlog in the database.
	Table = "create_ticket_logs"
)

// Columns holds all SQL columns for createticketlog fields.
var Columns = []string{
	FieldID,
	FieldUniqueID,
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

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"database-concurrency/ent/createticketlog"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
)

// CreateTicketLog is the model entity for the CreateTicketLog schema.
type CreateTicketLog struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// TicketID holds the value of the "ticket_id" field.
	TicketID uuid.UUID `json:"ticket_id,omitempty"`
	// UniqueID holds the value of the "unique_id" field.
	UniqueID uuid.UUID `json:"unique_id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

// scanValues returns the types for scanning values from sql.Rows.
func (*CreateTicketLog) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case createticketlog.FieldID:
			values[i] = new(sql.NullInt64)
		case createticketlog.FieldCreatedAt, createticketlog.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case createticketlog.FieldTicketID, createticketlog.FieldUniqueID:
			values[i] = new(uuid.UUID)
		default:
			return nil, fmt.Errorf("unexpected column %q for type CreateTicketLog", columns[i])
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the CreateTicketLog fields.
func (ctl *CreateTicketLog) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case createticketlog.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ctl.ID = int(value.Int64)
		case createticketlog.FieldTicketID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field ticket_id", values[i])
			} else if value != nil {
				ctl.TicketID = *value
			}
		case createticketlog.FieldUniqueID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field unique_id", values[i])
			} else if value != nil {
				ctl.UniqueID = *value
			}
		case createticketlog.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ctl.CreatedAt = value.Time
			}
		case createticketlog.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ctl.UpdatedAt = value.Time
			}
		}
	}
	return nil
}

// Update returns a builder for updating this CreateTicketLog.
// Note that you need to call CreateTicketLog.Unwrap() before calling this method if this CreateTicketLog
// was returned from a transaction, and the transaction was committed or rolled back.
func (ctl *CreateTicketLog) Update() *CreateTicketLogUpdateOne {
	return (&CreateTicketLogClient{config: ctl.config}).UpdateOne(ctl)
}

// Unwrap unwraps the CreateTicketLog entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ctl *CreateTicketLog) Unwrap() *CreateTicketLog {
	_tx, ok := ctl.config.driver.(*txDriver)
	if !ok {
		panic("ent: CreateTicketLog is not a transactional entity")
	}
	ctl.config.driver = _tx.drv
	return ctl
}

// String implements the fmt.Stringer.
func (ctl *CreateTicketLog) String() string {
	var builder strings.Builder
	builder.WriteString("CreateTicketLog(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ctl.ID))
	builder.WriteString("ticket_id=")
	builder.WriteString(fmt.Sprintf("%v", ctl.TicketID))
	builder.WriteString(", ")
	builder.WriteString("unique_id=")
	builder.WriteString(fmt.Sprintf("%v", ctl.UniqueID))
	builder.WriteString(", ")
	builder.WriteString("created_at=")
	builder.WriteString(ctl.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ctl.UpdatedAt.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// CreateTicketLogs is a parsable slice of CreateTicketLog.
type CreateTicketLogs []*CreateTicketLog

func (ctl CreateTicketLogs) config(cfg config) {
	for _i := range ctl {
		ctl[_i].config = cfg
	}
}
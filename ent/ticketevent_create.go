// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database-concurrency/ent/ticket"
	"database-concurrency/ent/ticketevent"
	"database-concurrency/ent/user"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TicketEventCreate is the builder for creating a TicketEvent entity.
type TicketEventCreate struct {
	config
	mutation *TicketEventMutation
	hooks    []Hook
}

// SetTicketID sets the "ticket_id" field.
func (tec *TicketEventCreate) SetTicketID(u uuid.UUID) *TicketEventCreate {
	tec.mutation.SetTicketID(u)
	return tec
}

// SetUserID sets the "user_id" field.
func (tec *TicketEventCreate) SetUserID(u uuid.UUID) *TicketEventCreate {
	tec.mutation.SetUserID(u)
	return tec
}

// SetType sets the "type" field.
func (tec *TicketEventCreate) SetType(s string) *TicketEventCreate {
	tec.mutation.SetType(s)
	return tec
}

// SetMetadata sets the "metadata" field.
func (tec *TicketEventCreate) SetMetadata(m map[string]interface{}) *TicketEventCreate {
	tec.mutation.SetMetadata(m)
	return tec
}

// SetVersions sets the "versions" field.
func (tec *TicketEventCreate) SetVersions(s string) *TicketEventCreate {
	tec.mutation.SetVersions(s)
	return tec
}

// SetNillableVersions sets the "versions" field if the given value is not nil.
func (tec *TicketEventCreate) SetNillableVersions(s *string) *TicketEventCreate {
	if s != nil {
		tec.SetVersions(*s)
	}
	return tec
}

// SetCreatedAt sets the "created_at" field.
func (tec *TicketEventCreate) SetCreatedAt(t time.Time) *TicketEventCreate {
	tec.mutation.SetCreatedAt(t)
	return tec
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (tec *TicketEventCreate) SetNillableCreatedAt(t *time.Time) *TicketEventCreate {
	if t != nil {
		tec.SetCreatedAt(*t)
	}
	return tec
}

// SetUpdatedAt sets the "updated_at" field.
func (tec *TicketEventCreate) SetUpdatedAt(t time.Time) *TicketEventCreate {
	tec.mutation.SetUpdatedAt(t)
	return tec
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (tec *TicketEventCreate) SetNillableUpdatedAt(t *time.Time) *TicketEventCreate {
	if t != nil {
		tec.SetUpdatedAt(*t)
	}
	return tec
}

// SetID sets the "id" field.
func (tec *TicketEventCreate) SetID(u uuid.UUID) *TicketEventCreate {
	tec.mutation.SetID(u)
	return tec
}

// SetNillableID sets the "id" field if the given value is not nil.
func (tec *TicketEventCreate) SetNillableID(u *uuid.UUID) *TicketEventCreate {
	if u != nil {
		tec.SetID(*u)
	}
	return tec
}

// SetUser sets the "user" edge to the User entity.
func (tec *TicketEventCreate) SetUser(u *User) *TicketEventCreate {
	return tec.SetUserID(u.ID)
}

// SetTicket sets the "ticket" edge to the Ticket entity.
func (tec *TicketEventCreate) SetTicket(t *Ticket) *TicketEventCreate {
	return tec.SetTicketID(t.ID)
}

// Mutation returns the TicketEventMutation object of the builder.
func (tec *TicketEventCreate) Mutation() *TicketEventMutation {
	return tec.mutation
}

// Save creates the TicketEvent in the database.
func (tec *TicketEventCreate) Save(ctx context.Context) (*TicketEvent, error) {
	var (
		err  error
		node *TicketEvent
	)
	tec.defaults()
	if len(tec.hooks) == 0 {
		if err = tec.check(); err != nil {
			return nil, err
		}
		node, err = tec.sqlSave(ctx)
	} else {
		var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
			mutation, ok := m.(*TicketEventMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}
			if err = tec.check(); err != nil {
				return nil, err
			}
			tec.mutation = mutation
			if node, err = tec.sqlSave(ctx); err != nil {
				return nil, err
			}
			mutation.id = &node.ID
			mutation.done = true
			return node, err
		})
		for i := len(tec.hooks) - 1; i >= 0; i-- {
			if tec.hooks[i] == nil {
				return nil, fmt.Errorf("ent: uninitialized hook (forgotten import ent/runtime?)")
			}
			mut = tec.hooks[i](mut)
		}
		v, err := mut.Mutate(ctx, tec.mutation)
		if err != nil {
			return nil, err
		}
		nv, ok := v.(*TicketEvent)
		if !ok {
			return nil, fmt.Errorf("unexpected node type %T returned from TicketEventMutation", v)
		}
		node = nv
	}
	return node, err
}

// SaveX calls Save and panics if Save returns an error.
func (tec *TicketEventCreate) SaveX(ctx context.Context) *TicketEvent {
	v, err := tec.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tec *TicketEventCreate) Exec(ctx context.Context) error {
	_, err := tec.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tec *TicketEventCreate) ExecX(ctx context.Context) {
	if err := tec.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tec *TicketEventCreate) defaults() {
	if _, ok := tec.mutation.CreatedAt(); !ok {
		v := ticketevent.DefaultCreatedAt()
		tec.mutation.SetCreatedAt(v)
	}
	if _, ok := tec.mutation.UpdatedAt(); !ok {
		v := ticketevent.DefaultUpdatedAt()
		tec.mutation.SetUpdatedAt(v)
	}
	if _, ok := tec.mutation.ID(); !ok {
		v := ticketevent.DefaultID()
		tec.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (tec *TicketEventCreate) check() error {
	if _, ok := tec.mutation.TicketID(); !ok {
		return &ValidationError{Name: "ticket_id", err: errors.New(`ent: missing required field "TicketEvent.ticket_id"`)}
	}
	if _, ok := tec.mutation.UserID(); !ok {
		return &ValidationError{Name: "user_id", err: errors.New(`ent: missing required field "TicketEvent.user_id"`)}
	}
	if _, ok := tec.mutation.GetType(); !ok {
		return &ValidationError{Name: "type", err: errors.New(`ent: missing required field "TicketEvent.type"`)}
	}
	if _, ok := tec.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "TicketEvent.created_at"`)}
	}
	if _, ok := tec.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "TicketEvent.updated_at"`)}
	}
	if _, ok := tec.mutation.UserID(); !ok {
		return &ValidationError{Name: "user", err: errors.New(`ent: missing required edge "TicketEvent.user"`)}
	}
	if _, ok := tec.mutation.TicketID(); !ok {
		return &ValidationError{Name: "ticket", err: errors.New(`ent: missing required edge "TicketEvent.ticket"`)}
	}
	return nil
}

func (tec *TicketEventCreate) sqlSave(ctx context.Context) (*TicketEvent, error) {
	_node, _spec := tec.createSpec()
	if err := sqlgraph.CreateNode(ctx, tec.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	return _node, nil
}

func (tec *TicketEventCreate) createSpec() (*TicketEvent, *sqlgraph.CreateSpec) {
	var (
		_node = &TicketEvent{config: tec.config}
		_spec = &sqlgraph.CreateSpec{
			Table: ticketevent.Table,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: ticketevent.FieldID,
			},
		}
	)
	if id, ok := tec.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := tec.mutation.GetType(); ok {
		_spec.SetField(ticketevent.FieldType, field.TypeString, value)
		_node.Type = value
	}
	if value, ok := tec.mutation.Metadata(); ok {
		_spec.SetField(ticketevent.FieldMetadata, field.TypeJSON, value)
		_node.Metadata = value
	}
	if value, ok := tec.mutation.Versions(); ok {
		_spec.SetField(ticketevent.FieldVersions, field.TypeString, value)
		_node.Versions = value
	}
	if value, ok := tec.mutation.CreatedAt(); ok {
		_spec.SetField(ticketevent.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := tec.mutation.UpdatedAt(); ok {
		_spec.SetField(ticketevent.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if nodes := tec.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   ticketevent.UserTable,
			Columns: []string{ticketevent.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: user.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.UserID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := tec.mutation.TicketIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   ticketevent.TicketTable,
			Columns: []string{ticketevent.TicketColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: ticket.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.TicketID = nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// TicketEventCreateBulk is the builder for creating many TicketEvent entities in bulk.
type TicketEventCreateBulk struct {
	config
	builders []*TicketEventCreate
}

// Save creates the TicketEvent entities in the database.
func (tecb *TicketEventCreateBulk) Save(ctx context.Context) ([]*TicketEvent, error) {
	specs := make([]*sqlgraph.CreateSpec, len(tecb.builders))
	nodes := make([]*TicketEvent, len(tecb.builders))
	mutators := make([]Mutator, len(tecb.builders))
	for i := range tecb.builders {
		func(i int, root context.Context) {
			builder := tecb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*TicketEventMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				nodes[i], specs[i] = builder.createSpec()
				var err error
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, tecb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, tecb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, tecb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (tecb *TicketEventCreateBulk) SaveX(ctx context.Context) []*TicketEvent {
	v, err := tecb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (tecb *TicketEventCreateBulk) Exec(ctx context.Context) error {
	_, err := tecb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tecb *TicketEventCreateBulk) ExecX(ctx context.Context) {
	if err := tecb.Exec(ctx); err != nil {
		panic(err)
	}
}

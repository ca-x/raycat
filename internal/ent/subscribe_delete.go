// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"raycat/internal/ent/predicate"
	"raycat/internal/ent/subscribe"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// SubscribeDelete is the builder for deleting a Subscribe entity.
type SubscribeDelete struct {
	config
	hooks    []Hook
	mutation *SubscribeMutation
}

// Where appends a list predicates to the SubscribeDelete builder.
func (sd *SubscribeDelete) Where(ps ...predicate.Subscribe) *SubscribeDelete {
	sd.mutation.Where(ps...)
	return sd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (sd *SubscribeDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, sd.sqlExec, sd.mutation, sd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (sd *SubscribeDelete) ExecX(ctx context.Context) int {
	n, err := sd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (sd *SubscribeDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(subscribe.Table, sqlgraph.NewFieldSpec(subscribe.FieldID, field.TypeUUID))
	if ps := sd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, sd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	sd.mutation.done = true
	return affected, err
}

// SubscribeDeleteOne is the builder for deleting a single Subscribe entity.
type SubscribeDeleteOne struct {
	sd *SubscribeDelete
}

// Where appends a list predicates to the SubscribeDelete builder.
func (sdo *SubscribeDeleteOne) Where(ps ...predicate.Subscribe) *SubscribeDeleteOne {
	sdo.sd.mutation.Where(ps...)
	return sdo
}

// Exec executes the deletion query.
func (sdo *SubscribeDeleteOne) Exec(ctx context.Context) error {
	n, err := sdo.sd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{subscribe.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (sdo *SubscribeDeleteOne) ExecX(ctx context.Context) {
	if err := sdo.Exec(ctx); err != nil {
		panic(err)
	}
}
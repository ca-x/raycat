// Code generated by ent, DO NOT EDIT.

package ent

import (
	"raycat/internal/ent/schema"
	"raycat/internal/ent/subscribe"
	"time"

	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	subscribeFields := schema.Subscribe{}.Fields()
	_ = subscribeFields
	// subscribeDescUpdateTimeoutSeconds is the schema descriptor for update_timeout_seconds field.
	subscribeDescUpdateTimeoutSeconds := subscribeFields[3].Descriptor()
	// subscribe.DefaultUpdateTimeoutSeconds holds the default value on creation for the update_timeout_seconds field.
	subscribe.DefaultUpdateTimeoutSeconds = subscribeDescUpdateTimeoutSeconds.Default.(int)
	// subscribeDescCreatedAt is the schema descriptor for created_at field.
	subscribeDescCreatedAt := subscribeFields[6].Descriptor()
	// subscribe.DefaultCreatedAt holds the default value on creation for the created_at field.
	subscribe.DefaultCreatedAt = subscribeDescCreatedAt.Default.(func() time.Time)
	// subscribeDescID is the schema descriptor for id field.
	subscribeDescID := subscribeFields[0].Descriptor()
	// subscribe.DefaultID holds the default value on creation for the id field.
	subscribe.DefaultID = subscribeDescID.Default.(func() uuid.UUID)
}

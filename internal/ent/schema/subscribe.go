package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

// Subscribe holds the schema definition for the Subscribe entity.
type Subscribe struct {
	ent.Schema
}

// Fields of the Subscribe.
func (Subscribe) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Enum("kind").Values("url_sub", "local_sub", "tpl_sub", "sub_entry", "none").Default("none"),
		field.String("location").Unique(),
		field.Int("update_timeout_seconds").Default(2),
		field.Int64("latency"),
		field.Time("expire_at"),
		field.Time("created_at").Default(time.Now),
	}
}

// Edges of the Subscribe.
func (Subscribe) Edges() []ent.Edge {
	return nil
}
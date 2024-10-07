// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// SubscribesColumns holds the columns for the "subscribes" table.
	SubscribesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUUID},
		{Name: "kind", Type: field.TypeEnum, Enums: []string{"url_sub", "local_sub", "tpl_sub", "sub_entry", "none"}, Default: "none"},
		{Name: "location", Type: field.TypeString, Unique: true},
		{Name: "update_timeout_seconds", Type: field.TypeInt, Default: 2},
		{Name: "latency", Type: field.TypeInt64},
		{Name: "expire_at", Type: field.TypeTime},
		{Name: "created_at", Type: field.TypeTime},
	}
	// SubscribesTable holds the schema information for the "subscribes" table.
	SubscribesTable = &schema.Table{
		Name:       "subscribes",
		Columns:    SubscribesColumns,
		PrimaryKey: []*schema.Column{SubscribesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		SubscribesTable,
	}
)

func init() {
}
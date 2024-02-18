package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("age").Positive(),
		field.String("name").Unique(),
		field.Float("cash").Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

// define indices
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// single index field
		index.Fields("name"),
		// composite index fields, make it unique to prevent duplicates
		index.Fields("age", "name").Unique(),
	}
}

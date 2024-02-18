package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Todo holds the schema definition for the Todo entity.
type Todo struct {
	ent.Schema
}

// Fields of the Todo.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.String("title"),
		field.Enum("status").Values("incomplete", "complete").Default("incomplete"),
	}
}

// Edges of the Todo.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		// Add an edge from Todo to User
		edge.From("user", User.Type).
			Ref("todos"). // This should match the name of the edge defined in the User schema
			Unique().     // Each Todo is linked to exactly one User
			Required(),   // (Optional) if every Todo must be associated with a User
	}
}

// Indexes of the Todo.
func (Todo) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("title").Unique(),
	}
}

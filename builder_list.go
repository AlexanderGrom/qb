package qb

import "strings"

// ListBuilder builds list of placeholders
type ListBuilder struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
}

// Append appends new values to the list
//  var b = new(qb.ListBuilder).Append("one", "two", "three").
//  _ = b.String() // $1, $2, $3
//  _ = b.Params() // ["one", "two", "three"]
func (b *ListBuilder) Append(values ...interface{}) *ListBuilder {
	if len(values) == 0 {
		return b
	}
	b.params = append(b.params, values...)
	b.groups = append(b.groups, func() string {
		return ", " + b.g().Placeholder(len(values))
	})
	return b
}

// String implementations Stringer interface
func (b *ListBuilder) String() string {
	if len(b.groups) == 0 {
		return ""
	}
	var s strings.Builder
	for _, f := range b.groups {
		s.WriteString(f())
	}
	return s.String()[2:]
}

// Params returns parameters for query
func (b *ListBuilder) Params() []interface{} {
	return b.params
}

// Grammar sets a Grammar
func (b *ListBuilder) Grammar(grammar Grammar) Builder {
	b.grammar = grammar
	return b
}

func (b *ListBuilder) g() Grammar {
	if b.grammar == nil {
		b.grammar = grammar()
	}
	return b.grammar
}

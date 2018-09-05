package qb

import (
	"strings"
)

// SetBuilder builds SET expressions
type SetBuilder struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
	regular bool
}

// Set adds a new SET expression
//  var b = new(qb.SetBuilder).Set("name", "Tom").Set("surname", "Johnson")
//  _ = b.String() // "name" = $1, "surname" = $2
//  _ = b.Params() // ["Tom", "Johnson"]
func (b *SetBuilder) Set(field string, value interface{}) *SetBuilder {
	b.params = append(b.params, value)
	b.groups = append(b.groups, func() string {
		return ", " + b.g().Wrap(field) + " = " + b.g().Placeholder(1)
	})
	return b
}

// SetRaw adds a new SET expression
//  var b = new(qb.SetBuilder).SetRaw("jsondata->'name' = %p", "Tom")
//  _ = b.String() // jsondata->'name' = $1
//  _ = b.Params() // ["Tom"]
func (b *SetBuilder) SetRaw(query string, params ...interface{}) *SetBuilder {
	var f = &format{
		query:   query,
		params:  params,
		grammar: b.g(),
	}
	b.params = append(b.params, f.Params()...)
	b.groups = append(b.groups, func() string {
		return ", " + f.Grammar(b.g()).String()
	})
	return b
}

// String implementations Stringer interface
func (b *SetBuilder) String() string {
	if len(b.groups) == 0 {
		return ""
	}
	defer b.r()
	var s strings.Builder
	for _, f := range b.groups {
		s.WriteString(f())
	}
	return s.String()[2:]
}

// Params returns parameters for query
func (b *SetBuilder) Params() []interface{} {
	return b.params
}

// Grammar sets a Grammar
func (b *SetBuilder) Grammar(grammar Grammar) Builder {
	b.grammar = grammar
	return b
}

func (b *SetBuilder) g() Grammar {
	if b.grammar == nil {
		b.grammar = grammar()
	}
	return b.grammar
}

func (b *SetBuilder) r() {
	if !b.regular {
		b.grammar = grammar()
	}
}

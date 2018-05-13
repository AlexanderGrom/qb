package qb

import (
	"strings"
)

// SetGroup groups SET expressions
type SetGroup struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
}

// Set adds a new SET expression
//  var b = new(qb.SetGroup).Set("name", "Tom").Set("surname", "Johnson")
//  _ = b.String() // "name" = $1, "surname" = $2
//  _ = b.Params() // ["Tom", "Johnson"]
func (g *SetGroup) Set(field string, value interface{}) *SetGroup {
	g.params = append(g.params, value)
	g.groups = append(g.groups, func() string {
		return ", " + g.g().Wrap(field) + " = " + g.g().Placeholder(1)
	})
	return g
}

// SetRaw adds a new SET expression
//  var b = new(qb.SetGroup).SetRaw("jsondata->'name' = %p", "Tom")
//  _ = b.String() // jsondata->'name' = $1
//  _ = b.Params() // ["Tom"]
func (g *SetGroup) SetRaw(query string, params ...interface{}) *SetGroup {
	var f = &format{
		query:   query,
		params:  params,
		grammar: g.g(),
	}
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return ", " + f.Grammar(g.g()).String()
	})
	return g
}

// String implementations Stringer interface
func (g *SetGroup) String() string {
	if len(g.groups) == 0 {
		return ""
	}
	var b strings.Builder
	for _, f := range g.groups {
		b.WriteString(f())
	}
	return b.String()[2:]
}

// Params returns parameters for query
func (g *SetGroup) Params() []interface{} {
	return g.params
}

// Grammar sets a Grammar
func (g *SetGroup) Grammar(grammar Grammar) Builder {
	g.grammar = grammar
	return g
}

func (g *SetGroup) g() Grammar {
	if g.grammar == nil {
		g.grammar = grammar()
	}
	return g.grammar
}

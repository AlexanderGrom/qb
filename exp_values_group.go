package qb

import "strings"

// ValuesGroup groups VALUES expressions
type ValuesGroup struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
}

// Values sets values and adds a new VALUES expression
//  var b = new(qb.ValuesGroup).
//    Values(1, "Marty", "McFly").
//    Values(2, "Emmett", "Brown")
//  _ = b.String() // ($1, $2, $3), ($4, $5, $6)
//  _ = b.Param()  // [1, "Marty", "McFly", 2, "Emmett", "Brown"]
func (g *ValuesGroup) Values(values ...interface{}) *ValuesGroup {
	g.params = append(g.params, values...)
	g.groups = append(g.groups, func() string {
		return ", (" + g.g().Placeholder(len(values)) + ")"
	})
	return g
}

// String implementations Stringer interface
func (g *ValuesGroup) String() string {
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
func (g *ValuesGroup) Params() []interface{} {
	return g.params
}

// Grammar sets a Grammar
func (g *ValuesGroup) Grammar(grammar Grammar) Builder {
	g.grammar = grammar
	return g
}

func (g *ValuesGroup) g() Grammar {
	if g.grammar == nil {
		g.grammar = grammar()
	}
	return g.grammar
}

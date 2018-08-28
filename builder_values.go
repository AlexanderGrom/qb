package qb

import "strings"

// ValuesBuilder builds VALUES expressions
type ValuesBuilder struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
}

// Values sets values and adds a new VALUES expression
//  var b = new(qb.ValuesBuilder).
//    Values(1, "Marty", "McFly").
//    Values(2, "Emmett", "Brown")
//  _ = b.String() // ($1, $2, $3), ($4, $5, $6)
//  _ = b.Params() // [1, "Marty", "McFly", 2, "Emmett", "Brown"]
func (b *ValuesBuilder) Values(values ...interface{}) *ValuesBuilder {
	b.params = append(b.params, values...)
	b.groups = append(b.groups, func() string {
		return ", (" + b.g().Placeholder(len(values)) + ")"
	})
	return b
}

// String implementations Stringer interface
func (b *ValuesBuilder) String() string {
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
func (b *ValuesBuilder) Params() []interface{} {
	return b.params
}

// Grammar sets a Grammar
func (b *ValuesBuilder) Grammar(grammar Grammar) Builder {
	b.grammar = grammar
	return b
}

func (b *ValuesBuilder) g() Grammar {
	if b.grammar == nil {
		b.grammar = grammar()
	}
	return b.grammar
}

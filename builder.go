// Package qb is a simple query builder.
//  var b = new(qb.WhereGroup).
//    Where("name", "=", "Marty").
//    Where("surname", "=", "McFly")
//
//  var q = qb.Query("SELECT id FROM table WHERE %s LIMIT %p", b, 1)
//  _ = b.String() // SELECT id FROM table WHERE "name" = $1 AND "surname" = $2 LIMIT $3
//  _ = b.Param()  // ["Marty", "McFly", 1]
package qb

import (
	"strings"
)

var (
	grammar  = PgsqlGrammar
	grammars = map[string]func() Grammar{}
)

type (
	// Grammar interface
	Grammar interface {
		Wrap(s string) string
		Placeholder(n int) string
	}

	// Builder interface
	Builder interface {
		String() string
		Params() []interface{}
		Grammar(Grammar) Builder
	}

	// Format query
	format struct {
		query   string
		params  []interface{}
		grammar Grammar
	}
)

// DefaultGrammar sets a default grammar
func DefaultGrammar(name string) {
	var ok bool
	if grammar, ok = grammars[name]; !ok {
		panic("qb: grammar '" + name + "' not found")
	}
}

// RegisterGrammar registers a new grammar
func RegisterGrammar(name string, grammar func() Grammar) {
	grammars[name] = grammar
}

// Query formats according to a format specifier and returns the sql query string
func Query(query string, params ...interface{}) Builder {
	return &format{
		query:   query,
		params:  params,
		grammar: grammar(),
	}
}

// String implementations Stringer interface
func (f *format) String() string {
	var (
		b strings.Builder
		p int
		s int
		r bool
	)
	for i := 0; i < len(f.query); i++ {
		switch {
		case f.query[i] == '%':
			if r = !r; !r {
				b.WriteString(f.query[s : i-1])
				b.WriteString(f.query[i : i+1])
				s = i + 1
			}
		case f.query[i] == 's' && r:
			if p >= len(f.params) {
				panic("qb: parameter not found")
			}
			if b, ok := f.params[p].(Builder); ok {
				b.Grammar(f.g())
			}
			b.WriteString(f.query[s : i-1])
			b.WriteString(toString(f.params[p]))
			s = i + 1
			p++
		case f.query[i] == 'p' && r:
			if p >= len(f.params) {
				panic("qb: parameter not found")
			}
			b.WriteString(f.query[s : i-1])
			b.WriteString(f.g().Placeholder(1))
			s = i + 1
			p++
		default:
			r = false
		}
	}
	b.WriteString(f.query[s:])
	return b.String()
}

// Params returns parameters for query
func (f *format) Params() []interface{} {
	var params = make([]interface{}, 0, len(f.params))
	for _, p := range f.params {
		if b, ok := p.(Builder); ok {
			params = append(params, b.Params()...)
		} else {
			params = append(params, p)
		}
	}
	return params
}

// Grammar sets a Grammar
func (f *format) Grammar(grammar Grammar) Builder {
	f.grammar = grammar
	return f
}

func (f *format) g() Grammar {
	if f.grammar == nil {
		f.grammar = grammar()
	}
	return f.grammar
}

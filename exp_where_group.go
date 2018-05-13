package qb

import "strings"

// WhereGroup groups WHERE expressions.
type WhereGroup struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
}

// Where adds an expression to the group
//  var b = new(qb.WhereGroup).Where("name", "=", "Tom")
//  _ = b.String() // "name" = $1
//  _ = b.Params() // ["Tom"]
func (g *WhereGroup) Where(field, operator string, value interface{}) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, value)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " " + operator + " " + g.g().Placeholder(1)
	})
	return g
}

// WhereOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereOr("id", "=", "1").WhereOr("id", "=", "2")
//  _ = b.String() // "id" = $1 OR "id" = $2
//  _ = b.Params() // [1, 2]
func (g *WhereGroup) WhereOr(field, operator string, value interface{}) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, value)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " " + operator + " " + g.g().Placeholder(1)
	})
	return g
}

// WhereRaw adds an expression to the group
//  var b = new(qb.WhereGroup).WhereRaw("jsondata->%p = %p", "name", "Tom")
//  _ = b.String() // jsondata->$1 = $2
//  _ = b.Params() // ["name", "Tom"]
func (g *WhereGroup) WhereRaw(query string, params ...interface{}) *WhereGroup {
	var (
		f = &format{
			query:  query,
			params: params,
		}
		b = g.and()
	)
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return b + f.Grammar(g.g()).String()
	})
	return g
}

// WhereRawOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereRawOr("jsondata->%p = %p", "name", "Tom")
//  _ = b.String() // jsondata->$1 = $2
//  _ = b.Params() // ["name", "Tom"]
func (g *WhereGroup) WhereRawOr(query string, params ...interface{}) *WhereGroup {
	var (
		f = &format{
			query:  query,
			params: params,
		}
		b = g.or()
	)
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return b + f.Grammar(g.g()).String()
	})
	return g
}

// WhereIn adds an expression to the group
//  var b = new(qb.WhereGroup).WhereIn("id", 1, 2, 3)
//  _ = b.String() // "id" IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (g *WhereGroup) WhereIn(field string, params ...interface{}) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IN (" + g.g().Placeholder(len(params)) + ")"
	})
	return g
}

// WhereInOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereInOr("id", 1, 2, 3)
//  _ = b.String() // "id" IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (g *WhereGroup) WhereInOr(field string, params ...interface{}) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IN (" + g.g().Placeholder(len(params)) + ")"
	})
	return g
}

// WhereNotIn adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotIn("id", 1, 2, 3)
//  _ = b.String() // "id" NOT IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (g *WhereGroup) WhereNotIn(field string, params ...interface{}) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " NOT IN (" + g.g().Placeholder(len(params)) + ")"
	})
	return g
}

// WhereNotInOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotInOr("id", 1, 2, 3)
//  _ = b.String() // "id" NOT IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (g *WhereGroup) WhereNotInOr(field string, params ...interface{}) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, params...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " NOT IN (" + g.g().Placeholder(len(params)) + ")"
	})
	return g
}

// WhereInSub adds an expression to the group
//  var b = new(qb.WhereGroup).WhereInSub("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (g *WhereGroup) WhereInSub(field string, query Builder) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, query.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IN (" + query.Grammar(g.g()).String() + ")"
	})
	return g
}

// WhereInSubOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereInSubOr("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (g *WhereGroup) WhereInSubOr(field string, query Builder) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, query.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IN (" + query.Grammar(g.g()).String() + ")"
	})
	return g
}

// WhereNotInSub adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotInSub("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" NOT IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (g *WhereGroup) WhereNotInSub(field string, query Builder) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, query.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " NOT IN (" + query.Grammar(g.g()).String() + ")"
	})
	return g
}

// WhereNotInSubOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotInSubOr("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" NOT IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (g *WhereGroup) WhereNotInSubOr(field string, query Builder) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, query.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " NOT IN (" + query.Grammar(g.g()).String() + ")"
	})
	return g
}

// WhereNull adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNull("data")
//  _ = b.String() // "data" IS NULL
//  _ = b.Params() // []
func (g *WhereGroup) WhereNull(field string) *WhereGroup {
	boolean := g.and()
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IS NULL"
	})
	return g
}

// WhereNullOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNullOr("data")
//  _ = b.String() // "data" IS NULL
//  _ = b.Params() // []
func (g *WhereGroup) WhereNullOr(field string) *WhereGroup {
	boolean := g.or()
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IS NULL"
	})
	return g
}

// WhereNotNull adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotNull("data")
//  _ = b.String() // "data" IS NOT NULL
//  _ = b.Params() // []
func (g *WhereGroup) WhereNotNull(field string) *WhereGroup {
	boolean := g.and()
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IS NOT NULL"
	})
	return g
}

// WhereNotNullOr adds an expression to the group
//  var b = new(qb.WhereGroup).WhereNotNullOr("data")
//  _ = b.String() // "data" IS NOT NULL
//  _ = b.Params() // []
func (g *WhereGroup) WhereNotNullOr(field string) *WhereGroup {
	boolean := g.or()
	g.groups = append(g.groups, func() string {
		return boolean + g.g().Wrap(field) + " IS NOT NULL"
	})
	return g
}

// WhereGroup adds an expression to the group
//  var g = new(qb.WhereGroup).Where("id", "=", 1).WhereOr("id", "=", 2)
//  var b = new(qb.WhereGroup).Where("name", "=", "Tom").WhereGroup(g)
//  _ = b.String() // "name" = $1 AND ("id" = $2 OR "id" = $3)
//  _ = b.Params() // ["Tom", 1, 2]
func (g *WhereGroup) WhereGroup(group *WhereGroup) *WhereGroup {
	boolean := g.and()
	g.params = append(g.params, group.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + "(" + group.Grammar(g.g()).String() + ")"
	})
	return g
}

// WhereGroupOr adds an expression to the group
//  var g = new(qb.WhereGroup).Where("id", "=", 1).WhereOr("id", "=", 2)
//  var b = new(qb.WhereGroup).Where("name", "=", "Tom").WhereGroupOr(g)
//  _ = b.String() // "name" = $1 AND ("id" = $2 OR "id" = $3)
//  _ = b.Params() // ["Tom", 1, 2]
func (g *WhereGroup) WhereGroupOr(group *WhereGroup) *WhereGroup {
	boolean := g.or()
	g.params = append(g.params, group.Params()...)
	g.groups = append(g.groups, func() string {
		return boolean + "(" + group.Grammar(g.g()).String() + ")"
	})
	return g
}

// String implementations Stringer interface
func (g *WhereGroup) String() string {
	var b strings.Builder
	for _, f := range g.groups {
		b.WriteString(f())
	}
	return b.String()
}

// Params returns parameters for query
func (g *WhereGroup) Params() []interface{} {
	return g.params
}

// Grammar sets a Grammar
func (g *WhereGroup) Grammar(grammar Grammar) Builder {
	g.grammar = grammar
	return g
}

func (g *WhereGroup) g() Grammar {
	if g.grammar == nil {
		g.grammar = grammar()
	}
	return g.grammar
}

func (g *WhereGroup) and() string {
	if len(g.groups) == 0 {
		return ""
	}
	return " AND "
}

func (g *WhereGroup) or() string {
	if len(g.groups) == 0 {
		return ""
	}
	return " OR "
}

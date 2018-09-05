package qb

import "strings"

// WhereBuilder builds WHERE expressions.
type WhereBuilder struct {
	groups  []func() string
	params  []interface{}
	grammar Grammar
	regular bool
}

// Where adds an expression to the group
//  var b = new(qb.WhereBuilder).Where("name", "=", "Tom")
//  _ = b.String() // "name" = $1
//  _ = b.Params() // ["Tom"]
func (b *WhereBuilder) Where(field, operator string, value interface{}) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, value)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " " + operator + " " + b.g().Placeholder(1)
	})
	return b
}

// WhereOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereOr("id", "=", "1").WhereOr("id", "=", "2")
//  _ = b.String() // "id" = $1 OR "id" = $2
//  _ = b.Params() // [1, 2]
func (b *WhereBuilder) WhereOr(field, operator string, value interface{}) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, value)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " " + operator + " " + b.g().Placeholder(1)
	})
	return b
}

// WhereRaw adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereRaw("jsondata->%p = %p", "name", "Tom")
//  _ = b.String() // jsondata->$1 = $2
//  _ = b.Params() // ["name", "Tom"]
func (b *WhereBuilder) WhereRaw(query string, params ...interface{}) *WhereBuilder {
	var (
		f = &format{
			query:  query,
			params: params,
		}
		s = b.and()
	)
	b.params = append(b.params, f.Params()...)
	b.groups = append(b.groups, func() string {
		return s + f.Grammar(b.g()).String()
	})
	return b
}

// WhereRawOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereRawOr("jsondata->%p = %p", "name", "Tom")
//  _ = b.String() // jsondata->$1 = $2
//  _ = b.Params() // ["name", "Tom"]
func (b *WhereBuilder) WhereRawOr(query string, params ...interface{}) *WhereBuilder {
	var (
		f = &format{
			query:  query,
			params: params,
		}
		s = b.or()
	)
	b.params = append(b.params, f.Params()...)
	b.groups = append(b.groups, func() string {
		return s + f.Grammar(b.g()).String()
	})
	return b
}

// WhereIn adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereIn("id", 1, 2, 3)
//  _ = b.String() // "id" IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (b *WhereBuilder) WhereIn(field string, params ...interface{}) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, params...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IN (" + b.g().Placeholder(len(params)) + ")"
	})
	return b
}

// WhereInOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereInOr("id", 1, 2, 3)
//  _ = b.String() // "id" IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (b *WhereBuilder) WhereInOr(field string, params ...interface{}) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, params...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IN (" + b.g().Placeholder(len(params)) + ")"
	})
	return b
}

// WhereNotIn adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotIn("id", 1, 2, 3)
//  _ = b.String() // "id" NOT IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (b *WhereBuilder) WhereNotIn(field string, params ...interface{}) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, params...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " NOT IN (" + b.g().Placeholder(len(params)) + ")"
	})
	return b
}

// WhereNotInOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotInOr("id", 1, 2, 3)
//  _ = b.String() // "id" NOT IN ($1, $2, $3)
//  _ = b.Params() // [1, 2, 3]
func (b *WhereBuilder) WhereNotInOr(field string, params ...interface{}) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, params...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " NOT IN (" + b.g().Placeholder(len(params)) + ")"
	})
	return b
}

// WhereInSub adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereInSub("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (b *WhereBuilder) WhereInSub(field string, query Builder) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, query.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IN (" + query.Grammar(b.g()).String() + ")"
	})
	return b
}

// WhereInSubOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereInSubOr("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (b *WhereBuilder) WhereInSubOr(field string, query Builder) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, query.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IN (" + query.Grammar(b.g()).String() + ")"
	})
	return b
}

// WhereNotInSub adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotInSub("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" NOT IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (b *WhereBuilder) WhereNotInSub(field string, query Builder) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, query.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " NOT IN (" + query.Grammar(b.g()).String() + ")"
	})
	return b
}

// WhereNotInSubOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotInSubOr("id", qb.Query(`SELECT id FROM table name = %p`, "Tom"))
//  _ = b.String() // "id" NOT IN (SELECT id FROM table name = $1)
//  _ = b.Params() // ["Tom"]
func (b *WhereBuilder) WhereNotInSubOr(field string, query Builder) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, query.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " NOT IN (" + query.Grammar(b.g()).String() + ")"
	})
	return b
}

// WhereNull adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNull("data")
//  _ = b.String() // "data" IS NULL
//  _ = b.Params() // []
func (b *WhereBuilder) WhereNull(field string) *WhereBuilder {
	boolean := b.and()
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IS NULL"
	})
	return b
}

// WhereNullOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNullOr("data")
//  _ = b.String() // "data" IS NULL
//  _ = b.Params() // []
func (b *WhereBuilder) WhereNullOr(field string) *WhereBuilder {
	boolean := b.or()
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IS NULL"
	})
	return b
}

// WhereNotNull adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotNull("data")
//  _ = b.String() // "data" IS NOT NULL
//  _ = b.Params() // []
func (b *WhereBuilder) WhereNotNull(field string) *WhereBuilder {
	boolean := b.and()
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IS NOT NULL"
	})
	return b
}

// WhereNotNullOr adds an expression to the group
//  var b = new(qb.WhereBuilder).WhereNotNullOr("data")
//  _ = b.String() // "data" IS NOT NULL
//  _ = b.Params() // []
func (b *WhereBuilder) WhereNotNullOr(field string) *WhereBuilder {
	boolean := b.or()
	b.groups = append(b.groups, func() string {
		return boolean + b.g().Wrap(field) + " IS NOT NULL"
	})
	return b
}

// WhereBuilder adds an expression to the group
//  var g = new(qb.WhereBuilder).Where("id", "=", 1).WhereOr("id", "=", 2)
//  var b = new(qb.WhereBuilder).Where("name", "=", "Tom").WhereBuilder(g)
//  _ = b.String() // "name" = $1 AND ("id" = $2 OR "id" = $3)
//  _ = b.Params() // ["Tom", 1, 2]
func (b *WhereBuilder) WhereBuilder(group *WhereBuilder) *WhereBuilder {
	boolean := b.and()
	b.params = append(b.params, group.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + "(" + group.Grammar(b.g()).String() + ")"
	})
	return b
}

// WhereBuilderOr adds an expression to the group
//  var g = new(qb.WhereBuilder).Where("id", "=", 1).WhereOr("id", "=", 2)
//  var b = new(qb.WhereBuilder).Where("name", "=", "Tom").WhereBuilderOr(g)
//  _ = b.String() // "name" = $1 AND ("id" = $2 OR "id" = $3)
//  _ = b.Params() // ["Tom", 1, 2]
func (b *WhereBuilder) WhereBuilderOr(group *WhereBuilder) *WhereBuilder {
	boolean := b.or()
	b.params = append(b.params, group.Params()...)
	b.groups = append(b.groups, func() string {
		return boolean + "(" + group.Grammar(b.g()).String() + ")"
	})
	return b
}

// String implementations Stringer interface
func (b *WhereBuilder) String() string {
	defer b.r()
	var s strings.Builder
	for _, f := range b.groups {
		s.WriteString(f())
	}
	return s.String()
}

// Params returns parameters for query
func (b *WhereBuilder) Params() []interface{} {
	return b.params
}

// Grammar sets a Grammar
func (b *WhereBuilder) Grammar(grammar Grammar) Builder {
	b.grammar = grammar
	b.regular = true
	return b
}

func (b *WhereBuilder) g() Grammar {
	if b.grammar == nil {
		b.grammar = grammar()
	}
	return b.grammar
}

func (b *WhereBuilder) r() {
	if !b.regular {
		b.grammar = grammar()
	}
}

func (b *WhereBuilder) and() string {
	if len(b.groups) == 0 {
		return ""
	}
	return " AND "
}

func (b *WhereBuilder) or() string {
	if len(b.groups) == 0 {
		return ""
	}
	return " OR "
}

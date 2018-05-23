package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_Query(t *testing.T) {
	b := new(WhereGroup).Where("name", "=", "test").WhereRaw("%s = %p", "field", "param")
	q := Query(
		"SELECT *, %%s, %%p, %, %%, %s%s FROM table WHERE status = %p AND %s ORDER BY %s %s LIMIT %p, OFFSET %p",
		"r2", "d2", "active", b, "id", "asc", 10, 0,
	)
	assert.Equal(t,
		`SELECT *, %s, %p, %, %, r2d2 FROM table WHERE status = $1 AND "name" = $2 AND field = $3 ORDER BY id asc LIMIT $4, OFFSET $5`,
		q.String(),
	)
	assert.Equal(t,
		[]interface{}{"active", "test", "param", 10, 0},
		q.Params(),
	)
}

func TestBuilder_DefaultGlammar(t *testing.T) {
	DefaultGrammar("mysql")
	b := new(WhereGroup).Where("name", "=", "test")
	q := Query(
		"SELECT *, %%s, %%p, %, %% FROM table WHERE status = %p AND %s LIMIT %p, OFFSET %p",
		"active", b, 10, 0,
	)
	assert.Equal(t,
		"SELECT *, %s, %p, %, % FROM table WHERE status = ? AND `name` = ? LIMIT ?, OFFSET ?",
		q.String(),
	)
	DefaultGrammar("postgres")
}

func TestBuilder_Glammar(t *testing.T) {
	b := new(WhereGroup).Where("name", "=", "test")
	q := Query(
		"SELECT *, %%s, %%p, %, %% FROM table WHERE status = %p AND %s LIMIT %p, OFFSET %p",
		"active", b, 10, 0,
	).Grammar(MysqlGrammar())
	assert.Equal(t,
		"SELECT *, %s, %p, %, % FROM table WHERE status = ? AND `name` = ? LIMIT ?, OFFSET ?",
		q.String(),
	)
	assert.Equal(t,
		[]interface{}{"active", "test", 10, 0},
		q.Params(),
	)
}

func BenchmarkBuilder_QueryString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b = new(WhereGroup).
			Where("type", "=", "a").
			WhereOr("type", "=", "b")

		_ = Query(`SELECT "name" FROM "table" WHERE %s LIMIT %p`, b, 10).String()
	}
}

func BenchmarkBuilder_QueryParams(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var b = new(WhereGroup).
			Where("type", "=", "a").
			WhereOr("type", "=", "b")

		_ = Query(`SELECT "name" FROM "table" WHERE %s LIMIT %p`, b, 10).Params()
	}
}

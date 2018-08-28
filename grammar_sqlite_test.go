package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSQLite_Wrap(t *testing.T) {
	var res string

	res = SQLiteGrammar().Wrap("name")
	assert.Equal(t, "`name`", res)

	res = SQLiteGrammar().Wrap("tx.name")
	assert.Equal(t, "`tx`.`name`", res)

	res = SQLiteGrammar().Wrap("public.tx.name")
	assert.Equal(t, "`public`.`tx`.`name`", res)
}

func TestSQLite_Placeholder(t *testing.T) {
	var res string

	res = SQLiteGrammar().Placeholder(0)
	assert.Equal(t, ``, res)

	res = SQLiteGrammar().Placeholder(1)
	assert.Equal(t, `?`, res)

	res = SQLiteGrammar().Placeholder(2)
	assert.Equal(t, `?, ?`, res)

	res = SQLiteGrammar().Placeholder(3)
	assert.Equal(t, `?, ?, ?`, res)
}

func BenchmarkSQLite_Wrap(b *testing.B) {
	var g = new(sqliteGrammar)
	for i := 0; i < b.N; i++ {
		_ = g.Wrap("public.test")
	}
}

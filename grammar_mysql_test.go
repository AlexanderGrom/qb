package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMySQL_Wrap(t *testing.T) {
	var res string
	res = MysqlGrammar().Wrap("name")
	assert.Equal(t, "`name`", res)

	res = MysqlGrammar().Wrap("tx.name")
	assert.Equal(t, "`tx`.`name`", res)

	res = MysqlGrammar().Wrap("public.tx.name")
	assert.Equal(t, "`public`.`tx`.`name`", res)
}

func TestMySQL_Placeholder(t *testing.T) {
	var res string

	res = MysqlGrammar().Placeholder(1)
	assert.Equal(t, `?`, res)

	res = MysqlGrammar().Placeholder(2)
	assert.Equal(t, `?, ?`, res)

	res = MysqlGrammar().Placeholder(3)
	assert.Equal(t, `?, ?, ?`, res)
}

func BenchmarkMySQL_Wrap(b *testing.B) {
	var g = new(mysqlGrammar)
	for i := 0; i < b.N; i++ {
		_ = g.Wrap("public.test")
	}
}

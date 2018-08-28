package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPgSQL_Wrap(t *testing.T) {
	var res string

	res = PgsqlGrammar().Wrap("name")
	assert.Equal(t, `"name"`, res)

	res = PgsqlGrammar().Wrap("tx.name")
	assert.Equal(t, `"tx"."name"`, res)

	res = PgsqlGrammar().Wrap("public.tx.name")
	assert.Equal(t, `"public"."tx"."name"`, res)

	res = PgsqlGrammar().Wrap("name::text")
	assert.Equal(t, `"name"::text`, res)

	res = PgsqlGrammar().Wrap("public.tx.name::text")
	assert.Equal(t, `"public"."tx"."name"::text`, res)

	res = PgsqlGrammar().Wrap("tx.name::")
	assert.Equal(t, `"tx"."name"::`, res)

	res = PgsqlGrammar().Wrap("tx.name:")
	assert.Equal(t, `"tx"."name":`, res)
}

func TestPgSQL_Placeholder(t *testing.T) {
	var res string

	res = PgsqlGrammar().Placeholder(0)
	assert.Equal(t, ``, res)

	res = PgsqlGrammar().Placeholder(1)
	assert.Equal(t, `$1`, res)

	res = PgsqlGrammar().Placeholder(2)
	assert.Equal(t, `$1, $2`, res)

	res = PgsqlGrammar().Placeholder(3)
	assert.Equal(t, `$1, $2, $3`, res)
}

func BenchmarkPgSQL_Wrap(b *testing.B) {
	var g = new(pgsqlGrammar)
	for i := 0; i < b.N; i++ {
		_ = g.Wrap("public.test::text")
	}
}

func BenchmarkPgSQL_Placeholder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = new(pgsqlGrammar).Placeholder(10)
	}
}

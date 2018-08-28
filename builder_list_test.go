package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	b := new(ListBuilder).Append("one", "two", "three")
	q := Query("SELECT id FROM table WHERE name ?| ARRAY[%s]", b)
	assert.Equal(t, `SELECT id FROM table WHERE name ?| ARRAY[$1, $2, $3]`, q.String())
	assert.Equal(t, []interface{}{"one", "two", "three"}, q.Params())
}

func TestList2(t *testing.T) {
	b := new(ListBuilder).
		Append("one", "two").
		Append("three").
		Append()
	q := Query("SELECT id FROM table WHERE name ?| ARRAY[%s]", b)
	assert.Equal(t, `SELECT id FROM table WHERE name ?| ARRAY[$1, $2, $3]`, q.String())
	assert.Equal(t, []interface{}{"one", "two", "three"}, q.Params())
}

func TestListMySQLGrammar(t *testing.T) {
	g := MysqlGrammar()
	b := new(ListBuilder).
		Append("one", "two").
		Append("three")
	q := Query("SELECT id FROM table WHERE name ?| ARRAY[%s]", b).Grammar(g)
	assert.Equal(t, `SELECT id FROM table WHERE name ?| ARRAY[?, ?, ?]`, q.String())
	assert.Equal(t, []interface{}{"one", "two", "three"}, q.Params())
}

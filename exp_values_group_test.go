package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	b := new(ValuesGroup).Values(1, "Marty", "McFly")
	q := Query("INSERT INTO table (id, name, surname) VALUES %s", b)
	assert.Equal(t, `INSERT INTO table (id, name, surname) VALUES ($1, $2, $3)`, q.String())
	assert.Equal(t, []interface{}{1, "Marty", "McFly"}, q.Params())
}

func TestValues2(t *testing.T) {
	b := new(ValuesGroup).
		Values(1, "Marty", "McFly").
		Values(2, "Emmett", "Brown")
	q := Query("INSERT INTO table (id, name, surname) VALUES %s", b)
	assert.Equal(t, `INSERT INTO table (id, name, surname) VALUES ($1, $2, $3), ($4, $5, $6)`, q.String())
	assert.Equal(t, []interface{}{1, "Marty", "McFly", 2, "Emmett", "Brown"}, q.Params())
}

func TestValuesMySQLGrammar(t *testing.T) {
	g := MysqlGrammar()
	b := new(ValuesGroup).
		Values(1, "Marty", "McFly").
		Values(2, "Emmett", "Brown")
	q := Query("INSERT INTO table (id, name, surname) VALUES %s", b).Grammar(g)
	assert.Equal(t, `INSERT INTO table (id, name, surname) VALUES (?, ?, ?), (?, ?, ?)`, q.String())
	assert.Equal(t, []interface{}{1, "Marty", "McFly", 2, "Emmett", "Brown"}, q.Params())
}

package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	b := new(SetGroup).
		Set("name", "Marty").
		Set("surname", "McFly")
	q := Query("UPDATE table SET %s WHERE id = %p", b, 10)
	assert.Equal(t, `UPDATE table SET "name" = $1, "surname" = $2 WHERE id = $3`, q.String())
	assert.Equal(t, []interface{}{"Marty", "McFly", 10}, q.Params())
}

func TestSetRaw(t *testing.T) {
	b := new(SetGroup).SetRaw("jsondata->%p = %p", "name", "Marty")
	q := Query("UPDATE table SET %s WHERE id = %p", b, 10)
	assert.Equal(t, `UPDATE table SET jsondata->$1 = $2 WHERE id = $3`, q.String())
	assert.Equal(t, []interface{}{"name", "Marty", 10}, q.Params())
}

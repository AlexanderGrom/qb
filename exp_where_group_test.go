package qb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWhere(t *testing.T) {
	b := new(WhereGroup).
		Where("type", "=", "a").
		WhereOr("type", "=", "b")

	assert.Equal(t, `"type" = $1 OR "type" = $2`, b.String())
	assert.Equal(t, []interface{}{"a", "b"}, b.Params())
}

func TestSelectWhere(t *testing.T) {
	b := new(WhereGroup).
		Where("type", "=", "a").
		WhereOr("type", "=", "b")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "type" = $2 OR "type" = $3 LIMIT $4`, q.String())
	assert.Equal(t, []interface{}{"param", "a", "b", 10}, q.Params())
}

func TestSelectWhereMySQLGrammar(t *testing.T) {
	g := MysqlGrammar()
	b := new(WhereGroup).
		Where("type", "=", "a").
		WhereOr("type", "=", "b")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10).Grammar(g)

	assert.Equal(t, "SELECT id FROM table WHERE param = ? AND `type` = ? OR `type` = ? LIMIT ?", q.String())
	assert.Equal(t, []interface{}{"param", "a", "b", 10}, q.Params())
}

func TestWhereRaw(t *testing.T) {
	b := new(WhereGroup).
		WhereRaw("class IN (%p, %p, %p)", 1, 2, 3).
		WhereRawOr("type IN (%p, %p, %p)", "a", "b", "c")

	assert.Equal(t, `class IN ($1, $2, $3) OR type IN ($4, $5, $6)`, b.String())
	assert.Equal(t, []interface{}{1, 2, 3, "a", "b", "c"}, b.Params())
}

func TestSelectWhereRaw(t *testing.T) {
	b := new(WhereGroup).
		WhereRaw("class IN (%p, %p, %p)", 1, 2, 3).
		WhereRawOr("type IN (%p, %p, %p)", "a", "b", "c")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND class IN ($2, $3, $4) OR type IN ($5, $6, $7) LIMIT $8`, q.String())
	assert.Equal(t, []interface{}{"param", 1, 2, 3, "a", "b", "c", 10}, q.Params())
}

func TestWhereIn(t *testing.T) {
	b := new(WhereGroup).
		WhereIn("class", 1, 2, 3).
		WhereInOr("type", "a", "b", "c")

	assert.Equal(t, `"class" IN ($1, $2, $3) OR "type" IN ($4, $5, $6)`, b.String())
	assert.Equal(t, []interface{}{1, 2, 3, "a", "b", "c"}, b.Params())
}

func TestSelectWhereIn(t *testing.T) {
	b := new(WhereGroup).
		WhereIn("class", 1, 2, 3).
		WhereInOr("type", "a", "b", "c")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "class" IN ($2, $3, $4) OR "type" IN ($5, $6, $7) LIMIT $8`, q.String())
	assert.Equal(t, []interface{}{"param", 1, 2, 3, "a", "b", "c", 10}, q.Params())
}

func TestWhereNotIn(t *testing.T) {
	b := new(WhereGroup).
		WhereNotIn("class", 1, 2, 3).
		WhereNotInOr("type", "a", "b", "c")

	assert.Equal(t, `"class" NOT IN ($1, $2, $3) OR "type" NOT IN ($4, $5, $6)`, b.String())
	assert.Equal(t, []interface{}{1, 2, 3, "a", "b", "c"}, b.Params())
}

func TestSelectWhereNotIn(t *testing.T) {
	b := new(WhereGroup).
		WhereNotIn("class", 1, 2, 3).
		WhereNotInOr("type", "a", "b", "c")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "class" NOT IN ($2, $3, $4) OR "type" NOT IN ($5, $6, $7) LIMIT $8`, q.String())
	assert.Equal(t, []interface{}{"param", 1, 2, 3, "a", "b", "c", 10}, q.Params())
}

func TestWhereInSub(t *testing.T) {
	b := new(WhereGroup).
		Where("status", "=", "active").
		WhereInSub("id", Query(`SELECT id FROM table name = %p`, "test")).
		WhereInSubOr("id", Query(`SELECT id FROM table surname = %p`, "best"))

	assert.Equal(t, `"status" = $1 AND "id" IN (SELECT id FROM table name = $2) OR "id" IN (SELECT id FROM table surname = $3)`, b.String())
	assert.Equal(t, []interface{}{"active", "test", "best"}, b.Params())
}

func TestSelectWhereInSub(t *testing.T) {
	b := new(WhereGroup).
		Where("status", "=", "active").
		WhereInSub("id", Query(`SELECT id FROM table name = %p`, "test")).
		WhereInSubOr("id", Query(`SELECT id FROM table surname = %p`, "best"))
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "status" = $2 AND "id" IN (SELECT id FROM table name = $3) OR "id" IN (SELECT id FROM table surname = $4) LIMIT $5`, q.String())
	assert.Equal(t, []interface{}{"param", "active", "test", "best", 10}, q.Params())
}

func TestWhereNotInSub(t *testing.T) {
	b := new(WhereGroup).
		Where("status", "=", "active").
		WhereNotInSub("id", Query(`SELECT id FROM table name = %p`, "test")).
		WhereNotInSubOr("id", Query(`SELECT id FROM table surname = %p`, "best"))

	assert.Equal(t, `"status" = $1 AND "id" NOT IN (SELECT id FROM table name = $2) OR "id" NOT IN (SELECT id FROM table surname = $3)`, b.String())
	assert.Equal(t, []interface{}{"active", "test", "best"}, b.Params())
}

func TestSelectWhereNotInSub(t *testing.T) {
	b := new(WhereGroup).
		Where("status", "=", "active").
		WhereNotInSub("id", Query(`SELECT id FROM table name = %p`, "test")).
		WhereNotInSubOr("id", Query(`SELECT id FROM table surname = %p`, "best"))
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "status" = $2 AND "id" NOT IN (SELECT id FROM table name = $3) OR "id" NOT IN (SELECT id FROM table surname = $4) LIMIT $5`, q.String())
	assert.Equal(t, []interface{}{"param", "active", "test", "best", 10}, q.Params())
}

func TestWhereNull(t *testing.T) {
	b := new(WhereGroup).
		WhereNull("updated_at").
		WhereNullOr("deleted_at")

	assert.Equal(t, `"updated_at" IS NULL OR "deleted_at" IS NULL`, b.String())
	assert.Nil(t, b.Params())
}

func TestSelectWhereNull(t *testing.T) {
	b := new(WhereGroup).
		WhereNull("updated_at").
		WhereNullOr("deleted_at")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "updated_at" IS NULL OR "deleted_at" IS NULL LIMIT $2`, q.String())
	assert.Equal(t, []interface{}{"param", 10}, q.Params())
}

func TestWhereNotNull(t *testing.T) {
	b := new(WhereGroup).
		WhereNotNull("updated_at").
		WhereNotNullOr("deleted_at")

	assert.Equal(t, `"updated_at" IS NOT NULL OR "deleted_at" IS NOT NULL`, b.String())
	assert.Nil(t, b.Params())
}

func TestSelectWhereNotNull(t *testing.T) {
	b := new(WhereGroup).
		WhereNotNull("updated_at").
		WhereNotNullOr("deleted_at")
	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND "updated_at" IS NOT NULL OR "deleted_at" IS NOT NULL LIMIT $2`, q.String())
	assert.Equal(t, []interface{}{"param", 10}, q.Params())
}

func TestWhereGroup(t *testing.T) {
	g1 := new(WhereGroup).
		Where("status", "=", "active").
		Where("type", "=", "a")

	g2 := new(WhereGroup).
		Where("status", "=", "passive").
		Where("type", "=", "b")

	b := new(WhereGroup).
		WhereGroup(g1).
		WhereGroupOr(g2)

	assert.Equal(t, `("status" = $1 AND "type" = $2) OR ("status" = $3 AND "type" = $4)`, b.String())
	assert.Equal(t, []interface{}{"active", "a", "passive", "b"}, b.Params())
}

func TestSelectWhereGroup(t *testing.T) {
	g1 := new(WhereGroup).
		Where("status", "=", "active").
		Where("type", "=", "a")

	g2 := new(WhereGroup).
		Where("status", "=", "passive").
		Where("type", "=", "b")

	b := new(WhereGroup).
		WhereGroup(g1).
		WhereGroupOr(g2)

	q := Query("SELECT id FROM table WHERE param = %p AND %s LIMIT %p", "param", b, 10)

	assert.Equal(t, `SELECT id FROM table WHERE param = $1 AND ("status" = $2 AND "type" = $3) OR ("status" = $4 AND "type" = $5) LIMIT $6`, q.String())
	assert.Equal(t, []interface{}{"param", "active", "a", "passive", "b", 10}, q.Params())
}

package qb

import (
	"unsafe"
)

type mysqlGrammar struct{}

var _ Grammar = (*mysqlGrammar)(nil)

func init() {
	RegisterGrammar("mysql", MysqlGrammar)
}

// MysqlGrammar returns a specific grammar for mysql
func MysqlGrammar() Grammar {
	return &mysqlGrammar{}
}

// Wrap wraps a string in quotes
func (g *mysqlGrammar) Wrap(s string) string {
	var dot int
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			dot++
		}
	}

	if dot == 0 {
		return "`" + s + "`"
	}

	var (
		n = dot*2 + 2 + len(s)
		b = make([]byte, n)
		w = 0
		x = 0
	)
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			w += copy(b[w:], "`"+s[x:i]+"`.")
			x = i + 1
		}
	}

	copy(b[w:], "`"+s[x:]+"`")
	return *(*string)(unsafe.Pointer(&b))
}

// Placeholder returns n count placeholders
func (g *mysqlGrammar) Placeholder(n int) string {
	if n < 0 {
		panic("qb: negative Placeholder count")
	}
	if n == 0 {
		return ""
	}
	if n == 1 {
		return "?"
	}

	var (
		p = ", ?"
		b = make([]byte, len(p)*n)
		w = copy(b, p)
	)

	for w < len(b) {
		copy(b[w:], b[:w])
		w *= 2
	}
	if len(b) >= len(p) {
		b = b[2:]
	}

	return *(*string)(unsafe.Pointer(&b))
}

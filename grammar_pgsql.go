package qb

import (
	"strconv"
	"unsafe"
)

type pgsqlGrammar struct {
	placeholders int
}

var _ Grammar = (*pgsqlGrammar)(nil)

func init() {
	RegisterGrammar("postgres", PgsqlGrammar)
}

// PgsqlGrammar returns a specific grammar for postgresql
func PgsqlGrammar() Grammar {
	return &pgsqlGrammar{}
}

// Wrap wraps a string in quotes
func (g *pgsqlGrammar) Wrap(s string) string {
	var (
		dot      int
		typecast bool
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			dot++
		case ':':
			typecast = true
			i = len(s)
			break
		}
	}

	if dot == 0 && !typecast {
		return `"` + s + `"`
	}

	var (
		n = dot*2 + 2 + len(s)
		b = make([]byte, n)
		w = 0
		x = 0
		t = false
	)
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '.':
			w += copy(b[w:], `"`+s[x:i]+`".`)
			x = i + 1
		case ':':
			w += copy(b[w:], `"`+s[x:i]+`":`)
			x = i + 1
			t = true
			i = len(s)
			break
		}
	}

	if !t {
		copy(b[w:], `"`+s[x:]+`"`)
	} else {
		copy(b[w:], s[x:])
	}

	return *(*string)(unsafe.Pointer(&b))
}

// Placeholder returns n count placeholders
func (g *pgsqlGrammar) placeholder() string {
	g.placeholders++
	return strconv.Itoa(g.placeholders)
}

func (g *pgsqlGrammar) Placeholder(n int) string {
	if n <= 0 {
		panic("qb: negative Placeholder count")
	}

	if n == 1 {
		return "$" + g.placeholder()
	}

	var (
		sep = ", "
		cap = len(sep)*(n-1) + n
	)
	for i := 1; i <= n; i++ {
		cap += intWeight(g.placeholders + i)
	}

	var (
		b = make([]byte, cap)
		w = copy(b, "$"+g.placeholder())
	)
	for i := 1; i < n; i++ {
		w += copy(b[w:], sep)
		w += copy(b[w:], "$"+g.placeholder())
	}

	return *(*string)(unsafe.Pointer(&b))
}

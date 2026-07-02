package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMapVo(t *testing.T) {
	t.Parallel()
	initial := map[string]int{"a": 1, "b": 2}
	m := NewMapVo(initial)
	assert.Equal(t, 2, m.Size())
	assert.True(t, m.Contains("a"))
	assert.Equal(t, 1, m.Get("a"))
}

func TestNewMapVo_NilInitial(t *testing.T) {
	t.Parallel()
	m := NewMapVo[string, int](nil)
	assert.Equal(t, 0, m.Size())
}

func TestNewEmptyMapVo(t *testing.T) {
	t.Parallel()
	m := NewEmptyMapVo[string, string]()
	assert.Equal(t, 0, m.Size())
	assert.False(t, m.Contains("key"))
}

func TestMapVoSetAndGet(t *testing.T) {
	t.Parallel()
	m := NewEmptyMapVo[string, int]()
	m.Set("x", 42)
	assert.True(t, m.Contains("x"))
	assert.Equal(t, 42, m.Get("x"))
}

func TestMapVoRemove(t *testing.T) {
	t.Parallel()
	m := NewMapVo(map[string]int{"a": 1, "b": 2})
	m.Remove("a")
	assert.False(t, m.Contains("a"))
	assert.True(t, m.Contains("b"))
	assert.Equal(t, 1, m.Size())
}

func TestMapVoValues(t *testing.T) {
	t.Parallel()
	m := NewMapVo(map[string]int{"k1": 10, "k2": 20})
	vals := m.Values()
	assert.Equal(t, 2, len(vals))
	assert.Equal(t, 10, vals["k1"])
	assert.Equal(t, 20, vals["k2"])
}

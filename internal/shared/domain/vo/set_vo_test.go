package vo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewSetVo(t *testing.T) {
	t.Parallel()
	s := NewSetVo("a", "b", "c")
	assert.Equal(t, 3, s.Size())
	assert.True(t, s.Contains("a"))
	assert.True(t, s.Contains("b"))
	assert.True(t, s.Contains("c"))
	assert.False(t, s.Contains("d"))
}

func TestNewEmptySetVo(t *testing.T) {
	t.Parallel()
	s := NewEmptySetVo[string]()
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Size())
}

func TestSetVoAdd(t *testing.T) {
	t.Parallel()
	s := NewEmptySetVo[int]()
	s.Add(1)
	s.Add(2)
	s.Add(1) // duplicate

	assert.Equal(t, 2, s.Size())
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2))
}

func TestSetVoRemove(t *testing.T) {
	t.Parallel()
	s := NewSetVo("x", "y")
	s.Remove("x")
	assert.False(t, s.Contains("x"))
	assert.True(t, s.Contains("y"))
	assert.Equal(t, 1, s.Size())
}

func TestSetVoIsEmpty(t *testing.T) {
	t.Parallel()
	s := NewEmptySetVo[string]()
	assert.True(t, s.IsEmpty())
	s.Add("item")
	assert.False(t, s.IsEmpty())
}

func TestSetVoValues(t *testing.T) {
	t.Parallel()
	s := NewSetVo(10, 20, 30)
	vals := s.Values()
	assert.Len(t, vals, 3)

	set := make(map[int]struct{})
	for _, v := range vals {
		set[v] = struct{}{}
	}
	assert.Contains(t, set, 10)
	assert.Contains(t, set, 20)
	assert.Contains(t, set, 30)
}

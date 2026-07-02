package ids

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateUuid(t *testing.T) {
	t.Parallel()
	u1 := GenerateUuid()
	u2 := GenerateUuid()

	assert.False(t, u1.IsEmpty(), "generated uuid should not be empty")
	assert.False(t, u2.IsEmpty(), "generated uuid should not be empty")
	assert.NotEqual(t, u1.String(), u2.String(), "two generated uuids should differ")
}

func TestNewUuid(t *testing.T) {
	t.Parallel()
	raw := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	u := NewUuid(raw)
	assert.Equal(t, "6ba7b810-9dad-11d1-80b4-00c04fd430c8", u.String())
}

func TestUuidEquals(t *testing.T) {
	t.Parallel()
	raw := uuid.MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	a := NewUuid(raw)
	b := NewUuid(raw)
	c := GenerateUuid()

	assert.True(t, a.Equals(b))
	assert.False(t, a.Equals(c))
}

func TestUuidIsEmpty(t *testing.T) {
	t.Parallel()
	empty := NewUuid(uuid.Nil)
	assert.True(t, empty.IsEmpty())

	nonEmpty := GenerateUuid()
	assert.False(t, nonEmpty.IsEmpty())
}

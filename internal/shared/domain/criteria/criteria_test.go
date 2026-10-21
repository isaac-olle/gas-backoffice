package criteria_test

import (
	"testing"

	"github.com/isaac-olle/gas-backoffice/internal/shared/domain/criteria"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCriteria_ToPostgresql(t *testing.T) {
	t.Parallel()

	t.Run("returns empty clause for no filters", func(t *testing.T) {
		t.Parallel()
		where, args, err := criteria.Criteria{}.ToPostgresql(1)
		require.NoError(t, err)
		assert.Empty(t, where)
		assert.Empty(t, args)
	})

	t.Run("builds a single condition starting at the given index", func(t *testing.T) {
		t.Parallel()
		c := criteria.Criteria{Filters: []criteria.Filter{
			{Field: "name", Operator: criteria.Equals, Value: "Acme"},
		}}
		where, args, err := c.ToPostgresql(1)
		require.NoError(t, err)
		assert.Equal(t, "name = $1", where)
		assert.Equal(t, []any{"Acme"}, args)
	})

	t.Run("joins multiple conditions with AND and offsets placeholders", func(t *testing.T) {
		t.Parallel()
		c := criteria.Criteria{Filters: []criteria.Filter{
			{Field: "name", Operator: criteria.Equals, Value: "Acme"},
			{Field: "user_id", Operator: criteria.Equals, Value: "u1"},
		}}
		where, args, err := c.ToPostgresql(2)
		require.NoError(t, err)
		assert.Equal(t, "name = $2 AND user_id = $3", where)
		assert.Equal(t, []any{"Acme", "u1"}, args)
	})

	t.Run("rejects field names that are not valid identifiers", func(t *testing.T) {
		t.Parallel()
		c := criteria.Criteria{Filters: []criteria.Filter{
			{Field: "name; DROP TABLE users;--", Operator: criteria.Equals, Value: "x"},
		}}
		_, _, err := c.ToPostgresql(1)
		assert.Error(t, err)
	})
}

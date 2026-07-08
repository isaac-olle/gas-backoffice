package event

import (
	sharedtypes "github.com/isaac-olle/gas-backoffice/internal/shared/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewEvent_Fields(t *testing.T) {
	t.Parallel()
	data := map[string]any{"foo": "bar"}
	before := time.Now()
	evt := NewEvent("agg-123", Type("user.created"), false, data, sharedtypes.AggregateAdmin)
	after := time.Now()

	assert.NotEmpty(t, evt.Id())
	assert.Equal(t, "agg-123", evt.AggregateID())
	assert.Equal(t, Type("user.created"), evt.Type())
	assert.False(t, evt.External())
	assert.Equal(t, sharedtypes.AggregateAdmin, evt.AggregateType())
	assert.Equal(t, data, evt.Data())
	assert.False(t, evt.OccurredOn().Before(before))
	assert.False(t, evt.OccurredOn().After(after))
}

func TestNewEvent_NilDataBecomesEmptyMap(t *testing.T) {
	t.Parallel()
	evt := NewEvent("agg-1", Type("x"), false, nil, sharedtypes.AggregateUser)
	assert.NotNil(t, evt.Data())
}

func TestNewEvent_External(t *testing.T) {
	t.Parallel()
	evt := NewEvent("agg-2", Type("ext.event"), true, nil, sharedtypes.AggregateClient)
	assert.True(t, evt.External())
}

func TestBaseEvent_Primitives(t *testing.T) {
	t.Parallel()
	data := map[string]any{"k": "v"}
	evt := NewEvent("agg-prim", Type("prim.event"), false, data, sharedtypes.AggregateScope)

	p := evt.Primitives()
	require.NotNil(t, p)
	assert.Equal(t, evt.Id(), p.ID)
	assert.Equal(t, "prim.event", p.Type)
	assert.Equal(t, "agg-prim", p.AggregateID)
	assert.Equal(t, sharedtypes.AggregateScope.String(), p.AggregateType)
	assert.Equal(t, data, p.Data)
	assert.False(t, p.External)
}

func TestAllTypes(t *testing.T) {
	t.Parallel()
	assert.Equal(t, Type("*"), AllTypes)
}

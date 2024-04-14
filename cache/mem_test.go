package cache

import (
	"context"
	"github.com/Electronic-Catalog/microkit/metric"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	mem, err := NewInMemoryCache(time.Second, WithMetricOption(metric.NewNop()))
	require.NoError(t, err)
	ctx, cn := context.WithCancel(context.Background())
	defer cn()

	err = mem.Set(ctx, "nop", "test1", "test1", time.Second*3)
	require.NoError(t, err)

	err = mem.Set(ctx, "nop", "test2", "test2", time.Second)
	require.NoError(t, err)

	time.Sleep(time.Second * 2)

	val, err := mem.GetKey(ctx, "nop", "test2")
	require.Error(t, NotFoundError, err)
	require.Equal(t, "", val)

	val, err = mem.GetKey(ctx, "nop", "test1")
	require.NoError(t, err)
	require.Equal(t, "test1", val)
}

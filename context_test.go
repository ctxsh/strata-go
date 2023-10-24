package strata

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContext(t *testing.T) {
	m := New(MetricsOpts{})
	ctx := IntoContext(context.Background(), m)
	m2, err := FromContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, m, m2)
}

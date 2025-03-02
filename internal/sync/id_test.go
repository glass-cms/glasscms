package sync_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/glass-cms/glasscms/internal/sync"
	"github.com/stretchr/testify/assert"
)

func TestSyncID_FromString(t *testing.T) {
	id := sync.ParseSyncID("sy_123e4567-e89b-12d3-a456-426614174000")
	assert.Equal(t, "sy_123e4567-e89b-12d3-a456-426614174000", id.String())
}

func TestSyncID_Intercept(t *testing.T) {
	id := sync.ParseSyncID("sy_123e4567-e89b-12d3-a456-426614174000")
	req, err := http.NewRequest(http.MethodGet, "http://example.com", nil)
	assert.NoError(t, err)
	assert.NoError(t, id.Intercept(context.Background(), req))
}

package sync

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

const (
	HeaderSyncID = "X-Sync-Id"
)

func NewSyncID() *ID {
	return &ID{id: uuid.New()}
}

func ParseSyncID(id string) *ID {
	if len(id) < 3 {
		return nil
	}
	return &ID{id: uuid.MustParse(id[3:])}
}

type ID struct {
	id uuid.UUID
}

func (s *ID) String() string {
	return fmt.Sprintf("sy_%s", s.id.String())
}

// Intercept will attach an X-Sync-Id header to the request
// and ensures that the sync ID is attached to the header.
func (s *ID) Intercept(_ context.Context, req *http.Request) error {
	req.Header.Set(HeaderSyncID, s.String())
	return nil
}

package sync_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/internal/sourcer/fs"
	"github.com/glass-cms/glasscms/internal/sync"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/log"
	"github.com/glass-cms/glasscms/pkg/mediatype"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ServerFunc func(w http.ResponseWriter, r *http.Request)

type TestServer struct {
	*httptest.Server

	itemListFunc      ServerFunc
	ItemListCallCount int
	ItemUpsertCalls   []api.ItemsUpsertJSONBody
}

func NewServer(listFunc ServerFunc) *TestServer {
	return &TestServer{itemListFunc: listFunc}
}

func stringPtr(s string) *string {
	return &s
}

func (s *TestServer) Init(t *testing.T) *TestServer {
	s.Server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/items" && r.Method == http.MethodGet {
			s.itemListFunc(w, r)
			s.ItemListCallCount++
			return
		}

		if r.URL.Path == "/items" && r.Method == http.MethodPatch {
			body, err := io.ReadAll(r.Body)
			assert.NoError(t, err)

			defer r.Body.Close()

			var req api.ItemsUpsertJSONBody
			err = json.Unmarshal(body, &req)
			assert.NoError(t, err)

			s.ItemUpsertCalls = append(s.ItemUpsertCalls, req)
			return
		}
	}))
	return s
}

func (s *TestServer) Close() {
	s.Server.Close()
}

//nolint:gocognit // Just deal with it.
func TestSyncer_Sync(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		livemode bool
		listFunc ServerFunc
		wantErr  error

		wantListCallCount   int
		wantUpsertCallCount int

		// Because we use the auto-gennerated Cobra docs as file system source, we can expect at least
		// x items to be upserted to the server, where x is the number of commands glasscms provides.
		wantNumberOfItemsUpsertedGte int
		upsertAssert                 func(t *testing.T, req api.ItemsUpsertJSONBody)
	}{
		`
		given a server with no items
		when the syncer is ran in livemode
		then the syncer should upsert at least 5 item to the server
		`: {
			livemode: true,
			listFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", mediatype.ApplicationJSON)
				w.WriteHeader(http.StatusOK)

				err := json.NewEncoder(w).Encode([]api.Item{})
				assert.NoError(t, err)
			},
			wantListCallCount:            1,
			wantUpsertCallCount:          1,
			wantNumberOfItemsUpsertedGte: 5,
		},
		`
		given a server with no items
		when the syncer is ran in livemode
		and the server returns an error
		then the syncer should return an error
		`: {
			livemode: true,
			listFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
			},
			wantErr:                      sync.ErrUnexpectedStatusCode,
			wantListCallCount:            1,
			wantUpsertCallCount:          0,
			wantNumberOfItemsUpsertedGte: 0,
		},
		`
		given a server with one item
		and the source does not contain the item
		when the syncer is ran in livemode
		the the syncer upserts the item with a delete flag
		`: {
			livemode: true,
			listFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", mediatype.ApplicationJSON)
				w.WriteHeader(http.StatusOK)

				err := json.NewEncoder(w).Encode([]api.Item{
					{
						Name:       "item",
						UpdateTime: time.Now(),
						Hash:       stringPtr("hash"),
					},
				})
				assert.NoError(t, err)
			},
			wantListCallCount:            1,
			wantUpsertCallCount:          1,
			wantNumberOfItemsUpsertedGte: 1,
			upsertAssert: func(t *testing.T, req api.ItemsUpsertJSONBody) {
				// find the item with the delete flag
				for _, item := range req {
					if item.Name == "item" {
						assert.NotNil(t, item.DeleteTime)
						return
					}
				}
				t.Fatal("item not found")
			},
		},
		`
		given a server with no items
		when the syncer is ran in drymode
		then the syncer should not upsert any items to the server
		`: {
			livemode: false,
			listFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", mediatype.ApplicationJSON)
				w.WriteHeader(http.StatusOK)

				err := json.NewEncoder(w).Encode([]api.Item{})
				assert.NoError(t, err)
			},
			wantListCallCount:            1,
			wantUpsertCallCount:          0,
			wantNumberOfItemsUpsertedGte: 0,
		},
		`
		given a server with one item
		and the source contains the same item with new content
		when the syncer is ran in livemode
		then the syncer should upsert the item with the new content
		`: {
			livemode: true,
			listFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", mediatype.ApplicationJSON)
				w.WriteHeader(http.StatusOK)

				err := json.NewEncoder(w).Encode([]api.Item{
					{
						Name:       "glasscms-completion",
						UpdateTime: time.Time{},
						Hash:       stringPtr("hash"),
					},
				})
				assert.NoError(t, err)
			},
			wantListCallCount:            1,
			wantUpsertCallCount:          1,
			wantNumberOfItemsUpsertedGte: 1,
			upsertAssert: func(t *testing.T, req api.ItemsUpsertJSONBody) {
				for _, item := range req {
					if item.Name == "glasscms-completion" {
						assert.NotEmpty(t, item.Content)
						return
					}
				}
				t.Fatal("item not found")
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			server := NewServer(tt.listFunc).Init(t)
			defer server.Close()

			client, err := api.NewClientWithResponses(server.URL)
			require.NoError(t, err)

			sourcer, err := fs.NewSourcer("../../docs/commands")
			require.NoError(t, err)

			syncer, err := sync.NewSyncer(sync.NewSyncID(), sourcer, client, log.NoopLogger(), &parser.Config{})
			require.NoError(t, err)

			// Act
			err = syncer.Sync(context.Background(), tt.livemode)
			if tt.wantErr != nil {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErr.Error())
				return
			}

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tt.wantListCallCount, server.ItemListCallCount)
			assert.Len(t, server.ItemUpsertCalls, tt.wantUpsertCallCount)

			// Assert the upserted items
			if tt.wantUpsertCallCount > 0 {
				req := server.ItemUpsertCalls[0]
				require.GreaterOrEqual(t, len(req), tt.wantNumberOfItemsUpsertedGte)

				if tt.upsertAssert != nil {
					tt.upsertAssert(t, req)
				}
			}
		})
	}
}

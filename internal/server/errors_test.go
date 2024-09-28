package server_test

import (
	"errors"
	"testing"

	"github.com/glass-cms/glasscms/internal/server"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/resource"
	"github.com/stretchr/testify/require"
)

func TestErrorMapperAlreadyExistsError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	tests := map[string]struct {
		args         args
		want         *api.Error
		expectPanics bool
	}{
		"maps resource.AlreadyExistsError to an API error response": {
			args: args{
				err: resource.NewAlreadyExistsError("item1", "item", errors.New("underlying error")),
			},
			want: &api.Error{
				Code:    api.ResourceAlreadyExists,
				Message: "An item with the name already exists",
				Type:    api.ApiError,
				Details: map[string]interface{}{
					"resource": "item",
					"name":     "item1",
				},
			},
		},
		"panics if error is not a resource.AlreadyExistsError": {
			args: args{
				err: errors.New("some error"),
			},
			expectPanics: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.expectPanics {
				require.Panics(t, func() {
					server.ErrorMapperAlreadyExistsError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, server.ErrorMapperAlreadyExistsError(tt.args.err))
		})
	}
}

func TestErrorMapperNotFoundError(t *testing.T) {
	t.Parallel()

	type args struct {
		err error
	}
	tests := map[string]struct {
		args         args
		want         *api.Error
		expectPanics bool
	}{
		"maps resource.NotFoundError to an API error response": {
			args: args{
				err: resource.NewNotFoundError("item1", "item", errors.New("underlying error")),
			},
			want: &api.Error{
				Code:    api.ResourceMissing,
				Message: "The item with the name was not found",
				Type:    api.ApiError,
				Details: map[string]interface{}{
					"resource": "item",
					"name":     "item1",
				},
			},
		},
		"panics if error is not a resource.NotFoundError": {
			args: args{
				err: errors.New("some error"),
			},
			expectPanics: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			if tt.expectPanics {
				require.Panics(t, func() {
					server.ErrorMapperNotFoundError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, server.ErrorMapperNotFoundError(tt.args.err))
		})
	}
}

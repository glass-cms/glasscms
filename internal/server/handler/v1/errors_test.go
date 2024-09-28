package v1_test

import (
	"errors"
	"testing"

	v1 "github.com/glass-cms/glasscms/api/v1"
	v1_handler "github.com/glass-cms/glasscms/internal/server/handler/v1"
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
		want         *v1.Error
		expectPanics bool
	}{
		"maps resource.AlreadyExistsError to an API error response": {
			args: args{
				err: resource.NewAlreadyExistsError("item1", "item", errors.New("underlying error")),
			},
			want: &v1.Error{
				Code:    v1.ResourceAlreadyExists,
				Message: "An item with the name already exists",
				Type:    v1.ApiError,
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
					v1_handler.ErrorMapperAlreadyExistsError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, v1_handler.ErrorMapperAlreadyExistsError(tt.args.err))
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
		want         *v1.Error
		expectPanics bool
	}{
		"maps resource.NotFoundError to an API error response": {
			args: args{
				err: resource.NewNotFoundError("item1", "item", errors.New("underlying error")),
			},
			want: &v1.Error{
				Code:    v1.ResourceMissing,
				Message: "The item with the name was not found",
				Type:    v1.ApiError,
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
					v1_handler.ErrorMapperNotFoundError(tt.args.err)
				})
				return
			}

			require.Equal(t, tt.want, v1_handler.ErrorMapperNotFoundError(tt.args.err))
		})
	}
}

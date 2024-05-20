package item_test

import (
	"testing"

	"github.com/glass-cms/glasscms/item"
	"github.com/stretchr/testify/assert"
)

func TestTitle(t *testing.T) {
	t.Run("when properties is nil", func(t *testing.T) {
		t.Parallel()
		i := &item.Item{}
		title := i.Title()
		assert.Nil(t, title)
	})

	t.Run("when properties does not contain title", func(t *testing.T) {
		t.Parallel()
		i := &item.Item{
			Properties: map[string]interface{}{
				"notTitle": "value",
			},
		}
		title := i.Title()
		assert.Nil(t, title)
	})

	t.Run("when properties contains title but not string", func(t *testing.T) {
		t.Parallel()
		i := &item.Item{
			Properties: map[string]interface{}{
				"title": 123,
			},
		}
		title := i.Title()
		assert.Nil(t, title)
	})

	t.Run("when properties contains title and is string", func(t *testing.T) {
		t.Parallel()
		expectedTitle := "expectedTitle"
		i := &item.Item{
			Properties: map[string]interface{}{
				"title": expectedTitle,
			},
		}
		title := i.Title()
		assert.NotNil(t, title)
		assert.Equal(t, expectedTitle, *title)
	})
}

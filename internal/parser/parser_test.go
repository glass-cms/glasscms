package parser_test

import (
	"strings"
	"testing"
	"time"

	"github.com/glass-cms/glasscms/internal/parser"
	"github.com/glass-cms/glasscms/pkg/wikilink"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockSource struct {
	name   string
	reader *strings.Reader
}

func NewMockSource(name string, data string) *MockSource {
	return &MockSource{
		name:   name,
		reader: strings.NewReader(data),
	}
}

func (m *MockSource) Read(p []byte) (int, error) {
	return m.reader.Read(p)
}

func (m *MockSource) Close() error {
	return nil
}

func (m *MockSource) Name() string {
	return m.name
}

func (m *MockSource) CreateTime() time.Time {
	return time.Now()
}

func (m *MockSource) UpdateTime() time.Time {
	return time.Now()
}

func TestParse(t *testing.T) {
	// Arrange
	source := NewMockSource("test", "---\ntitle: Test\n---\n# Test\n")

	// Act
	item, err := parser.Parse(source)

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, "test", item.Name)
	assert.Equal(t, "\n# Test\n", item.Content)
	assert.Equal(t, map[string]interface{}{"title": "Test"}, item.Properties)
}

func TestParseWithHiddenProperty(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		content        string
		hiddenProperty string
		hiddenValue    bool
		shouldBeHidden bool
	}{
		{
			name:           "Hidden with truthy value and truthy config",
			content:        "---\ntitle: Test\ndraft: true\n---\n# Test\n",
			hiddenProperty: "draft",
			hiddenValue:    true,
			shouldBeHidden: true,
		},
		{
			name:           "Not hidden with falsy value and truthy config",
			content:        "---\ntitle: Test\ndraft: false\n---\n# Test\n",
			hiddenProperty: "draft",
			hiddenValue:    true,
			shouldBeHidden: false,
		},
		{
			name:           "Hidden with falsy value and falsy config",
			content:        "---\ntitle: Test\ndraft: false\n---\n# Test\n",
			hiddenProperty: "draft",
			hiddenValue:    false,
			shouldBeHidden: true,
		},
		{
			name:           "Not hidden with truthy value and falsy config",
			content:        "---\ntitle: Test\ndraft: true\n---\n# Test\n",
			hiddenProperty: "draft",
			hiddenValue:    false,
			shouldBeHidden: false,
		},
		{
			name:           "Hidden with string truthy value",
			content:        "---\ntitle: Test\nstatus: draft\n---\n# Test\n",
			hiddenProperty: "status",
			hiddenValue:    true,
			shouldBeHidden: false, // "draft" is not considered truthy
		},
		{
			name:           "Not hidden when property doesn't exist",
			content:        "---\ntitle: Test\n---\n# Test\n",
			hiddenProperty: "draft",
			hiddenValue:    true,
			shouldBeHidden: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Arrange
			source := NewMockSource("test", tc.content)
			config := parser.Config{
				HiddenProperty: tc.hiddenProperty,
				HiddenValue:    tc.hiddenValue,
			}

			// Act
			item, err := parser.ParseWithConfig(source, config)

			// Assert
			if tc.shouldBeHidden {
				assert.Nil(t, item, "Item should be nil when hidden")
				assert.ErrorIs(t, err, parser.ErrItemHidden, "Error should be ErrItemHidden when item is hidden")
			} else {
				assert.NotNil(t, item, "Item should not be nil when not hidden")
				assert.NoError(t, err, "Error should be nil when item is not hidden")
			}
		})
	}
}

func TestParseWikilinks(t *testing.T) {
	t.Parallel()

	source := NewMockSource("test", "---\ntitle: Test\n---\n# Test [[link1]] [[link2|Link Two]]")
	item, err := parser.Parse(source)

	require.NoError(t, err)
	// Assert that there is a wikilinks property in metadata
	assert.NotNil(t, item.Metadata["wikilinks"])
	// Assert that there are two wikilinks
	links, ok := item.Metadata["wikilinks"].([]wikilink.Link)
	require.True(t, ok)
	assert.Len(t, links, 2)

	// Assert that the wikilinks are correct
	assert.Equal(t, "link1", links[0].Target)
	assert.Equal(t, "link2", links[1].Target)
	assert.Equal(t, "Link Two", links[1].DisplayText)
	assert.Equal(t, "[[link2|Link Two]]", links[1].Original)
}

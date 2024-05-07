package parser_test

import (
	"testing"
	"time"

	"strings"

	"github.com/glass-cms/glasscms/parser"
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

func (m *MockSource) CreatedAt() time.Time {
	return time.Now()
}

func (m *MockSource) ModifiedAt() time.Time {
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
	assert.Equal(t, "test", item.Title)
	assert.Equal(t, "# Test\n", item.Content)
	assert.Equal(t, map[string]interface{}{"title": "Test"}, item.Properties)
}

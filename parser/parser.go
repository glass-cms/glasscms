package parser

import (
	"bytes"
	"errors"
	"io"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/sourcer"
	"gopkg.in/yaml.v3"
)

const (
	numParts = 3
)

// Parse reads the content of a source and returns an item.
func Parse(src sourcer.Source) (*item.Item, error) {
	c, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	// Split the content into front matter and markdown
	parts := bytes.SplitN(c, []byte("---\n"), numParts)
	if len(parts) < numParts {
		return nil, errors.New("invalid content")
	}

	// Parse the YAML front matter
	var properties map[string]any
	err = yaml.Unmarshal(parts[1], &properties)
	if err != nil {
		return nil, err
	}

	// Keep the markdown content as is
	content := string(parts[2])

	return &item.Item{
		Title:      src.Name(),
		Content:    content,
		CreateTime: src.CreatedAt(),
		UpdateTime: src.ModifiedAt(),
		Properties: properties,
	}, nil
}

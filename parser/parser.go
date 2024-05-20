package parser

import (
	"bytes"
	"errors"
	"io"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/sourcer"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidFrontMatter = errors.New("invalid front matter yaml")
)

const (
	seperatorBytes = 4
)

func Parse(src sourcer.Source) (*item.Item, error) {
	c, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	frontMatter, content, err := extractFrontMatter(c)
	if err != nil {
		return nil, err
	}

	var properties map[string]any
	if len(frontMatter) > 0 {
		err = yaml.Unmarshal(frontMatter, &properties)
		if err != nil {
			return nil, err
		}
	}

	return &item.Item{
		Name:       src.Name(),
		Content:    string(content),
		CreateTime: src.CreatedAt(),
		UpdateTime: src.ModifiedAt(),
		Properties: properties,
	}, nil
}

func extractFrontMatter(content []byte) ([]byte, []byte, error) {
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, content, nil
	}

	frontMatterEnd := bytes.Index(content[seperatorBytes:], []byte("\n---\n"))
	if frontMatterEnd == -1 {
		return nil, nil, ErrInvalidFrontMatter
	}

	frontMatterEnd += seperatorBytes // Account for the initial "---\n"
	frontMatter := content[seperatorBytes:frontMatterEnd]
	markdown := content[frontMatterEnd+seperatorBytes:]

	return frontMatter, markdown, nil
}

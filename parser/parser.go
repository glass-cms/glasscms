package parser

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/glass-cms/glasscms/internal/item"
	"github.com/glass-cms/glasscms/sourcer"
	"github.com/mozillazg/go-slugify"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidFrontMatter = errors.New("invalid front matter yaml")
)

const (
	seperatorBytes = 4
)

// Parse reads the content of a source and extracts the front matter and markdown content.
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

	pathname := src.Name()

	// TODO: I want the name to be slugified but the slashes to be preserved.

	return &item.Item{
		Name:        slugify.Slugify(pathname),
		DisplayName: nameFromPath(pathname),
		Content:     string(content),
		Hash:        HashContent(content),
		CreateTime:  src.CreatedAt(),
		UpdateTime:  src.ModifiedAt(),
		Properties:  properties,
	}, nil
}

// HashContent generates a hash for the content.
func HashContent(content []byte) string {
	hasher := sha256.New()
	hasher.Write(content)
	return hex.EncodeToString(hasher.Sum(nil))
}

func nameFromPath(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	return name
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

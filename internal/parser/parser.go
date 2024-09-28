package parser

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
	"strings"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/slug"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidFrontMatter = errors.New("invalid front matter yaml")
)

const (
	seperatorBytes = 4
)

// Parse reads the content of a source and extracts the front matter and markdown content.
func Parse(src sourcer.Source) (*api.Item, error) {
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
	name := slug.Slug(pathname, slug.AllowSlashesOption())

	return &api.Item{
		Name:        name,
		DisplayName: nameFromPath(pathname),
		Content:     string(content),
		CreateTime:  src.CreateTime(),
		UpdateTime:  src.UpdateTime(),
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

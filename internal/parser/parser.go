package parser

import (
	"bytes"
	"errors"
	"io"
	"maps"
	"path/filepath"
	"strings"

	"github.com/glass-cms/glasscms/internal/sourcer"
	"github.com/glass-cms/glasscms/pkg/api"
	"github.com/glass-cms/glasscms/pkg/slug"
	"github.com/glass-cms/glasscms/pkg/wikilink"
	"gopkg.in/yaml.v3"
)

var (
	ErrInvalidFrontMatter = errors.New("invalid front matter yaml")
	// ErrItemHidden is returned when an item is hidden based on the hidden property configuration.
	ErrItemHidden = errors.New("item is hidden based on front matter property")
)

const (
	seperatorBytes = 4
)

// Config holds configuration options for the parser.
type Config struct {
	// HiddenProperty is the name of the front matter property that determines if an item is hidden.
	// If empty, no filtering based on visibility is performed.
	HiddenProperty string

	// HiddenValue determines how to interpret the hidden property value.
	// If true, truthy values (true, "true", "yes", "1", etc.) indicate the item is hidden.
	// If false, falsy values (false, "false", "no", "0", etc.) indicate the item is hidden.
	HiddenValue bool

	// ParseWikilinks determines if wikilinks should be parsed.
	ParseWikilinks bool

	// AdditionalMetadata is a map of additional metadata to be added to the item.
	AdditionalMetadata map[string]any
}

// Parse reads the content of a source and extracts the front matter and markdown content.
func Parse(src sourcer.Source) (*api.Item, error) {
	return ParseWithConfig(src, Config{
		ParseWikilinks: true,
	})
}

// ParseWithConfig reads the content of a source and extracts the front matter and markdown content,
// applying the provided configuration options.
func ParseWithConfig(src sourcer.Source, config Config) (*api.Item, error) {
	c, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	defer src.Close()

	frontMatter, content, err := extractFrontMatter(c)
	if err != nil {
		return nil, err
	}
	contentStr := string(content)

	var properties map[string]any
	if len(frontMatter) > 0 {
		err = yaml.Unmarshal(frontMatter, &properties)
		if err != nil {
			return nil, err
		}
	}

	// Check if the item should be considered hidden based on the config
	if config.HiddenProperty != "" && properties != nil {
		if propValue, exists := properties[config.HiddenProperty]; exists {
			isHidden := isTruthy(propValue)

			if (!config.HiddenValue && !isHidden) || (config.HiddenValue && isHidden) {
				return nil, ErrItemHidden
			}
		}
	}

	metadata := make(map[string]any)
	if config.ParseWikilinks {
		links := wikilink.ParseLinks(contentStr)
		if len(links) > 0 {
			metadata["wikilinks"] = links
		}
	}

	if config.AdditionalMetadata != nil {
		maps.Copy(metadata, config.AdditionalMetadata)
	}

	pathname := src.Name()
	name := slug.Slug(pathname, slug.AllowSlashesOption())

	hash, err := api.HashItem(contentStr, properties, metadata)
	if err != nil {
		return nil, err
	}

	return &api.Item{
		Name:        name,
		DisplayName: nameFromPath(pathname),
		Content:     contentStr,
		Hash:        &hash,
		CreateTime:  src.CreateTime(),
		UpdateTime:  src.UpdateTime(),
		Properties:  properties,
		Metadata:    metadata,
	}, nil
}

func isTruthy(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case string:
		lv := strings.ToLower(v)
		return lv == "true" || lv == "yes" || lv == "1" || lv == "on"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v != 0
	case float32, float64:
		return v != 0
	default:
		return false
	}
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

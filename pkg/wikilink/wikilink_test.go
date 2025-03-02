package wikilink_test

import (
	"testing"

	"github.com/glass-cms/glasscms/pkg/wikilink"
	"github.com/stretchr/testify/assert"
)

func TestParseLinks(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		content  string
		expected []wikilink.Link
	}{
		{
			name:    "Simple WikiLink",
			content: "This is a [[simple-link]] in text.",
			expected: []wikilink.Link{
				{
					Target:      "simple-link",
					DisplayText: "simple-link",
					Original:    "[[simple-link]]",
				},
			},
		},
		{
			name:    "WikiLink with display text",
			content: "This is a [[target|Display Text]] in text.",
			expected: []wikilink.Link{
				{
					Target:      "target",
					DisplayText: "Display Text",
					Original:    "[[target|Display Text]]",
				},
			},
		},
		{
			name:    "Multiple wikilink",
			content: "This has [[link1]] and [[link2|Link Two]] in it.",
			expected: []wikilink.Link{
				{
					Target:      "link1",
					DisplayText: "link1",
					Original:    "[[link1]]",
				},
				{
					Target:      "link2",
					DisplayText: "Link Two",
					Original:    "[[link2|Link Two]]",
				},
			},
		},
		{
			name:     "No wikilink",
			content:  "This has no wiki links in it.",
			expected: []wikilink.Link{},
		},
		{
			name:    "WikiLink with spaces",
			content: "This is a [[page with spaces]] in text.",
			expected: []wikilink.Link{
				{
					Target:      "page with spaces",
					DisplayText: "page with spaces",
					Original:    "[[page with spaces]]",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			links := wikilink.ParseLinks(tt.content)
			assert.Equal(t, tt.expected, links)
		})
	}
}

// Package wikilinks provides functionality for parsing and processing WikiLinks in markdown content.
package wikilink

import (
	"regexp"
	"strings"
)

var (
	// WikiLinkRegex matches WikiLinks in the format [[link-target]] or [[link-target|link-text]].
	WikiLinkRegex = regexp.MustCompile(`\[\[([^|\]]+)(?:\|([^\]]+))?\]\]`)
)

// Link represents a WikiLink with its target and optional display text.
type Link struct {
	Target      string `json:"target"`
	DisplayText string `json:"display_text"`
	Original    string `json:"original"`
}

// ParseLinks extracts all WikiLinks from the given markdown content.
func ParseLinks(content string) []Link {
	matches := WikiLinkRegex.FindAllStringSubmatch(content, -1)
	links := make([]Link, 0, len(matches))

	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		original := match[0]
		target := strings.TrimSpace(match[1])
		displayText := target

		// If there's a custom display text (format: [[target|display]]).
		if len(match) > 2 && match[2] != "" {
			displayText = strings.TrimSpace(match[2])
		}

		links = append(links, Link{
			Target:      target,
			DisplayText: displayText,
			Original:    original,
		})
	}

	return links
}

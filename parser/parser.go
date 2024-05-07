package parser

import (
	"io"

	"github.com/glass-cms/glasscms/item"
	"github.com/glass-cms/glasscms/sourcer"
)

type Parser struct {
	Config Config
}

func NewParser(config Config) *Parser {
	return &Parser{
		Config: config,
	}
}

type Config struct {
}

func (p *Parser) Parse(src sourcer.Source) (*item.Item, error) {
	_, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	// TODO: Parse the yaml front matter and markdown content from the source.

	return &item.Item{}, nil
}

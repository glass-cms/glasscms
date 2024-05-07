package parser

import "github.com/glass-cms/glasscms/item"

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

func (p *Parser) Parse(_ string) (*item.Item, error) {
	return &item.Item{}, nil
}

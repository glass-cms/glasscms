package parser

import "github.com/glass-cms/glasscms/item"

type Parser struct {
	Config ParserConfig
}

func NewParser(config ParserConfig) *Parser {
	return &Parser{
		Config: config,
	}
}

type ParserConfig struct {
}

func (p *Parser) Parse(data string) (*item.Item, error) {
	return nil, nil
}

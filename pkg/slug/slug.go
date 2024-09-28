package slug

import (
	"regexp"
	"strings"

	"github.com/rainycape/unidecode"
)

const (
	defaultSeparator       = "-"
	validCharRe            = `[^a-zA-Z0-9]+`
	validCharWithSlashesRe = `[^a-zA-Z0-9/]+`
)

type Option func(*slugOptions)

type slugOptions struct {
	AllowSlashes bool
	Separator    string
}

func AllowSlashesOption() func(*slugOptions) {
	return func(opts *slugOptions) {
		opts.AllowSlashes = true
	}
}

func CustomSeparatorOption(separator string) func(*slugOptions) {
	return func(opts *slugOptions) {
		opts.Separator = separator
	}
}

func Slug(s string, opts ...Option) string {
	s = unidecode.Unidecode(s)

	options := &slugOptions{
		Separator: defaultSeparator,
	}

	for _, opt := range opts {
		opt(options)
	}

	var re *regexp.Regexp
	if options.AllowSlashes {
		re = regexp.MustCompile(validCharWithSlashesRe)
	} else {
		re = regexp.MustCompile(validCharRe)
	}
	s = re.ReplaceAllString(s, options.Separator)

	s = strings.Trim(s, options.Separator)
	re = regexp.MustCompile(regexp.QuoteMeta(options.Separator) + `+`)
	s = re.ReplaceAllString(s, options.Separator)
	s = strings.ToLower(s)
	return s
}

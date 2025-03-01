package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glass-cms/glasscms/pkg/slug"
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

const (
	DocsCommandsFolder = "./docs/commands"
)

type DocsCommand struct {
	Command *cobra.Command

	linkPrefix   string
	slugifyLinks bool
}

// NewDocsCommand creates a new cobra.Command for `docs` which generates
// documentation for the application.
func NewDocsCommand() *DocsCommand {
	dc := &DocsCommand{}
	dc.Command = &cobra.Command{
		Use:    "docs",
		Short:  "Generate documentation",
		RunE:   dc.Execute,
		Hidden: true, // Development commands are hidden.
		Args:   cobra.NoArgs,
	}

	dc.Command.Flags().StringVar(
		&dc.linkPrefix,
		"link-prefix",
		"",
		"Prefix to add to documentation links (e.g., '../commands/')",
	)
	dc.Command.Flags().BoolVar(
		&dc.slugifyLinks,
		"slugify-links",
		false,
		"Convert link names to slugs",
	)

	return dc
}

func (c *DocsCommand) Execute(_ *cobra.Command, _ []string) error {
	if _, err := os.Stat(DocsCommandsFolder); os.IsNotExist(err) {
		err = os.MkdirAll(DocsCommandsFolder, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	linkHandlerFunc := func(name string) string {
		if c.slugifyLinks {
			name = slug.Slug(strings.TrimSuffix(name, filepath.Ext(name)))
		}
		return c.linkPrefix + name
	}

	return doc.GenMarkdownTreeCustom(rootCmd, DocsCommandsFolder, DocFilePrepender, linkHandlerFunc)
}

func DocFilePrepender(filename string) string {
	type FrontMatter struct {
		Title           string `yaml:"title"`
		CreateTimestamp int64  `yaml:"createTime"`
	}

	name := filepath.Base(filename)
	name = strings.TrimSuffix(name, filepath.Ext(name))

	title := strings.ReplaceAll(name, "_", " ")
	title = strings.ReplaceAll(title, "-", " ")
	title = cases.Title(language.English).String(title)

	frontMatter := FrontMatter{
		Title:           title,
		CreateTimestamp: time.Now().Unix(),
	}

	yamlFrontMatter, err := yaml.Marshal(&frontMatter)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	return "---\n" + string(yamlFrontMatter) + "---\n"
}

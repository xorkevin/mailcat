package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var (
	docOutputDir string
)

var (
	// docCmd represents the doc command
	docCmd = &cobra.Command{
		Use:               "doc",
		Short:             "generate documentation for mailcat",
		Long:              `generate documentation for mailcat in several formats`,
		DisableAutoGenTag: true,
	}

	docManCmd = &cobra.Command{
		Use:   "man",
		Short: "generate man page documentation for mailcat",
		Long:  `generate man page documentation for mailcat`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := doc.GenManTree(rootCmd, &doc.GenManHeader{
				Title:   "mailcat",
				Section: "1",
			}, docOutputDir); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
		DisableAutoGenTag: true,
	}

	docMdCmd = &cobra.Command{
		Use:   "md",
		Short: "generate markdown documentation for mailcat",
		Long:  `generate markdown documentation for mailcat`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := doc.GenMarkdownTree(rootCmd, docOutputDir); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		},
		DisableAutoGenTag: true,
	}
)

func init() {
	rootCmd.AddCommand(docCmd)
	docCmd.AddCommand(docManCmd)
	docCmd.AddCommand(docMdCmd)

	docCmd.PersistentFlags().StringVarP(&docOutputDir, "output", "o", ".", "documentation output path")
}

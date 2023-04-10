package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type (
	Cmd struct {
		rootCmd     *cobra.Command
		version     string
		rootFlags   rootFlags
		formatFlags formatFlags
		sendFlags   sendFlags
		docFlags    docFlags
	}

	rootFlags struct {
		logLevel string
	}
)

func New() *Cmd {
	return &Cmd{}
}

func (c *Cmd) Execute() {
	buildinfo := ReadVCSBuildInfo()
	c.version = buildinfo.ModVersion
	rootCmd := &cobra.Command{
		Use:               "mailcat",
		Short:             "A mail and smtp test tool",
		Long:              `A mail and smtp test tool`,
		Version:           c.version,
		DisableAutoGenTag: true,
	}
	c.rootCmd = rootCmd

	rootCmd.AddCommand(c.getFormatCmd())
	rootCmd.AddCommand(c.getSendCmd())
	rootCmd.AddCommand(c.getDocCmd())

	if err := rootCmd.Execute(); err != nil {
		c.logFatal(err)
		return
	}
}

func (c *Cmd) logFatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

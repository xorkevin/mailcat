package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/formatter"
)

var (
	formatInteractive bool
	formatTmpdir      string
	formatCRLFOutput  bool
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Formats plaintext mail output",
	Long:  `Formats plaintext mail output`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := formatter.Format(os.Stdin, os.Stdout, formatter.Opts{
			CRLF: formatCRLFOutput,
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.PersistentFlags().BoolVarP(&formatInteractive, "interactive", "i", false, "specify mail headers interactively")
	formatCmd.PersistentFlags().StringVarP(&formatTmpdir, "tmpdir", "d", "", "tmpdir for mail ($TMPDIR/mailcat if unset)")
	formatCmd.PersistentFlags().BoolVarP(&formatCRLFOutput, "crlf", "b", false, "output with CRLF line endings")
}

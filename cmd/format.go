package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/formatter"
)

var (
	formatInteractive bool
	formatTmpdir      string
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Formats plaintext mail output",
	Long:  `Formats plaintext mail output`,
	Run: func(cmd *cobra.Command, args []string) {
		formatter.Format(os.Stdin)
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.PersistentFlags().BoolVarP(&formatInteractive, "interactive", "i", false, "specify mail headers interactively")
	formatCmd.PersistentFlags().StringVarP(&formatTmpdir, "tmpdir", "d", "", "tmpdir for mail ($TMPDIR/mailcat if unset)")
}

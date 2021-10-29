package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	formatInteractive bool
	formatTmpdir      string
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Prints component configs",
	Long:  `Prints component configs`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mail format")
	},
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.PersistentFlags().BoolVarP(&formatInteractive, "interactive", "i", false, "specify mail headers interactively")
	formatCmd.PersistentFlags().StringVarP(&formatTmpdir, "tmpdir", "d", "", "tmpdir for mail ($TMPDIR/mailcat if unset)")
}

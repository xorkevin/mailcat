package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/send"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends smtp mail",
	Long:  `Sends smtp mail`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := send.Send(os.Stdin, send.Opts{}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(sendCmd)
}

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
	formatBodyInput   bool
	formatHeaders     []string
	formatAddHeaders  []string
	formatMsgIDDomain string
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Formats plaintext mail output",
	Long:  `Formats plaintext mail output`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := formatter.Format(os.Stdin, os.Stdout, formatter.Opts{
			CRLF:        formatCRLFOutput,
			Body:        formatBodyInput,
			Headers:     formatHeaders,
			AddHeaders:  formatAddHeaders,
			MsgIDDomain: formatMsgIDDomain,
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.PersistentFlags().StringVarP(&formatTmpdir, "tmpdir", "d", "", "tmpdir for mail ($TMPDIR/mailcat if unset)")
	formatCmd.PersistentFlags().BoolVarP(&formatCRLFOutput, "crlf", "m", false, "output with CRLF line endings")
	formatCmd.PersistentFlags().BoolVarP(&formatBodyInput, "body", "b", false, "input is body instead of a full RFC5322 message with headers")
	formatCmd.PersistentFlags().StringArrayVarP(&formatHeaders, "header", "s", nil, "set default header value (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringArrayVarP(&formatAddHeaders, "add", "a", nil, "specify header values to be added (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringVarP(&formatMsgIDDomain, "msgid", "y", "mail.example.com", "set default generated message id domain")
}

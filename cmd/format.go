package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/formatter"
)

var (
	formatCRLFOutput  bool
	formatBodyInput   bool
	formatHeaders     []string
	formatAddHeaders  []string
	formatMsgIDDomain string
	formatEdit        bool
	formatEmpty       bool
)

var formatCmd = &cobra.Command{
	Use:   "fmt",
	Short: "Formats plaintext mail output",
	Long:  `Formats plaintext mail output`,
	Run: func(cmd *cobra.Command, args []string) {
		var r io.Reader = os.Stdin
		if formatEmpty {
			r = bytes.NewReader(nil)
		}
		if err := formatter.Format(r, os.Stdout, formatter.Opts{
			CRLF:        formatCRLFOutput,
			Body:        formatBodyInput,
			Headers:     formatHeaders,
			AddHeaders:  formatAddHeaders,
			MsgIDDomain: formatMsgIDDomain,
			Edit:        formatEdit,
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(formatCmd)

	formatCmd.PersistentFlags().BoolVarP(&formatCRLFOutput, "crlf", "m", false, "output with CRLF line endings")
	formatCmd.PersistentFlags().BoolVarP(&formatBodyInput, "body", "b", false, "input is body instead of a full RFC5322 message with headers")
	formatCmd.PersistentFlags().StringArrayVarP(&formatHeaders, "header", "s", nil, "set default header value (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringArrayVarP(&formatAddHeaders, "add", "a", nil, "specify header values to be added (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringVarP(&formatMsgIDDomain, "msgid", "y", "mail.example.com", "set default generated message id domain")
	formatCmd.PersistentFlags().BoolVarP(&formatEdit, "edit", "e", false, "output in editor convenient format")
	formatCmd.PersistentFlags().BoolVarP(&formatEmpty, "empty", "z", false, "do not read from stdin and instead use empty reader")
}

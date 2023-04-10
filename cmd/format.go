package cmd

import (
	"bytes"
	"io"
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/formatter"
)

type (
	formatFlags struct {
		opts  formatter.Opts
		empty bool
	}
)

func (c *Cmd) getFormatCmd() *cobra.Command {
	formatCmd := &cobra.Command{
		Use:               "fmt",
		Short:             "Formats plaintext mail output",
		Long:              `Formats plaintext mail output`,
		Run:               c.execFormatCmd,
		DisableAutoGenTag: true,
	}
	formatCmd.PersistentFlags().BoolVarP(&c.formatFlags.opts.CRLF, "crlf", "m", false, "output with CRLF line endings")
	formatCmd.PersistentFlags().BoolVarP(&c.formatFlags.opts.Body, "body", "b", false, "input is body instead of a full RFC5322 message with headers")
	formatCmd.PersistentFlags().StringArrayVarP(&c.formatFlags.opts.Headers, "header", "s", nil, "set default header value (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringArrayVarP(&c.formatFlags.opts.AddHeaders, "add", "a", nil, "specify header values to be added (HEADER:VALUE); may be specified multiple times")
	formatCmd.PersistentFlags().StringVarP(&c.formatFlags.opts.MsgIDDomain, "msgid", "y", "mail.example.com", "set default generated message id domain")
	formatCmd.PersistentFlags().BoolVarP(&c.formatFlags.opts.Edit, "edit", "e", false, "output in editor convenient format")
	formatCmd.PersistentFlags().BoolVarP(&c.formatFlags.empty, "empty", "z", false, "do not read from stdin and instead use empty reader")

	return formatCmd
}

func (c *Cmd) execFormatCmd(cmd *cobra.Command, args []string) {
	var r io.Reader = os.Stdin
	if c.formatFlags.empty {
		r = bytes.NewReader(nil)
	}
	if err := formatter.Format(r, os.Stdout, c.formatFlags.opts); err != nil {
		c.logFatal(err)
		return
	}
}

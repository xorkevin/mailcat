package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/send"
)

type (
	sendFlags struct {
		opts             send.Opts
		sendAddr         string
		sendUsername     string
		sendPassword     string
		sendFrom         string
		sendTo           string
		sendDKIMSelector string
		sendDKIMKeyFile  string
	}
)

func (c *Cmd) getSendCmd() *cobra.Command {
	sendCmd := &cobra.Command{
		Use:               "send",
		Short:             "Sends smtp mail",
		Long:              `Sends smtp mail`,
		Run:               c.execSendCmd,
		DisableAutoGenTag: true,
	}
	sendCmd.PersistentFlags().StringVarP(&c.sendFlags.opts.Addr, "server", "s", "", "smtp server address")
	sendCmd.PersistentFlags().StringVarP(&c.sendFlags.opts.Username, "username", "u", "", "smtp auth username")
	sendCmd.PersistentFlags().StringVarP(&c.sendFlags.opts.Password, "password", "a", "", "smtp auth password")
	sendCmd.PersistentFlags().StringVarP(&c.sendFlags.opts.From, "from", "i", "", "smtp from")
	sendCmd.PersistentFlags().StringVarP(&c.sendFlags.opts.To, "to", "o", "", "smtp to")
	sendCmd.PersistentFlags().StringVar(&c.sendFlags.opts.DKIMSelector, "dkim-selector", "", "dkim selector")
	sendCmd.PersistentFlags().StringVar(&c.sendFlags.opts.DKIMKeyFile, "dkim-keyfile", "", "dkim key file (PEM)")
	return sendCmd
}

func (c *Cmd) execSendCmd(cmd *cobra.Command, args []string) {
	if err := send.Send(os.Stdin, c.sendFlags.opts); err != nil {
		c.logFatal(err)
		return
	}
}

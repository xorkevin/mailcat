package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"xorkevin.dev/mailcat/send"
)

var (
	sendAddr         string
	sendUsername     string
	sendPassword     string
	sendFrom         string
	sendTo           string
	sendDKIMSelector string
	sendDKIMKeyFile  string
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "Sends smtp mail",
	Long:  `Sends smtp mail`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := send.Send(os.Stdin, send.Opts{
			Addr:         sendAddr,
			Username:     sendUsername,
			Password:     sendUsername,
			From:         sendFrom,
			To:           sendTo,
			DKIMSelector: sendDKIMSelector,
			DKIMKeyFile:  sendDKIMKeyFile,
		}); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	},
	DisableAutoGenTag: true,
}

func init() {
	rootCmd.AddCommand(sendCmd)

	sendCmd.PersistentFlags().StringVarP(&sendAddr, "server", "s", "", "smtp server address")
	sendCmd.PersistentFlags().StringVarP(&sendUsername, "username", "u", "", "smtp auth username")
	sendCmd.PersistentFlags().StringVarP(&sendPassword, "password", "a", "", "smtp auth password")
	sendCmd.PersistentFlags().StringVarP(&sendFrom, "from", "i", "", "smtp from")
	sendCmd.PersistentFlags().StringVarP(&sendTo, "to", "o", "", "smtp to")
	sendCmd.PersistentFlags().StringVar(&sendDKIMSelector, "dkim-selector", "", "dkim selector")
	sendCmd.PersistentFlags().StringVar(&sendDKIMKeyFile, "dkim-keyfile", "", "dkim key file (PEM)")
}

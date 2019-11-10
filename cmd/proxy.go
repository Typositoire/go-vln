package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/typositoire/go-vln/proxy"
)

// proxyCmd represents the api command
var proxyCmd = &cobra.Command{
	Use:   "proxy",
	Short: "A brief description of proxy command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
TODO`,
	Run: func(cmd *cobra.Command, args []string) {
		p, err := proxy.NewProxyClient()

		if err != nil {
			panic(err)
		}

		p.Run()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringP("backend", "b", "vault", "SymlinkDB Backend")
	proxyCmd.Flags().StringP("vault-addr", "A", "http://127.0.0.1:8200", "Address of the vault server where symlinkDB is.")

	viper.BindPFlags(proxyCmd.Flags())
}

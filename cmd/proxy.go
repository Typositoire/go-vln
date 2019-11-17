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
		options := make(map[string]string)
		options["backend"] = viper.GetString("backend")
		options["hostURL"] = viper.GetString("vault-addr")
		options["beVaultAddr"] = viper.GetString("be-vault-addr")
		options["beVaultSymlinkDBPath"] = viper.GetString("be-vault-symlinkdbpath")
		options["beFilePath"] = viper.GetString("be-file-path")
		options["beGitRepository"] = viper.GetString("be-git-repository")
		options["beGitAccessToken"] = viper.GetString("be-git-accesstoken")

		p, err := proxy.NewProxyClient(options)

		if err != nil {
			panic(err)
		}

		p.Run()
	},
}

func init() {
	rootCmd.AddCommand(proxyCmd)

	proxyCmd.Flags().StringP("vault-addr", "A", "http://127.0.0.1:8200", "Address of the vault server where symlinkDB is.")
	proxyCmd.Flags().StringP("vault-app-role-id", "r", "placeHolder", "Vault app role id for symlinksDB access.")
	proxyCmd.Flags().StringP("vault-app-role-secret", "R", "placeHolder", "Vault app role secret for symlinksDB access.")

	proxyCmd.Flags().StringP("backend", "b", "vault", "SymlinkDB backend")

	// FileBackend
	proxyCmd.Flags().StringP("be-file-path", "f", "./db.json", "Path to file containing the KV Database in json.")

	// VaultBackend
	proxyCmd.Flags().StringP("be-vault-addr", "a", "http://127.0.0.1:8200", "Vault ADDR to get secrets from.")
	proxyCmd.Flags().StringP("be-vault-symlinkdbpath", "p", "secret/data/vln/symlinksdb", "Path to secret containing symlinks KV.")

	// GitBackend
	proxyCmd.Flags().StringP("be-git-accesstoken", "t", "access_token", "Git access token to pull the DB.")
	proxyCmd.Flags().StringP("be-git-repository", "T", "https://github.com/Typositoire/test-db", "Git repository which contains db.json at root.")

	viper.BindPFlags(proxyCmd.Flags())
}

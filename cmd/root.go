package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	symlinkDBPath string
	cfgFile       string
	logLevel      string
	backend       string
	filePath      string
	vaultAddr     string
	version       string
)

var rootCmd = &cobra.Command{
	Use:   "go-vln",
	Short: "go-vln is a vault symlink proxy",
	Long: `A longer text about
	
	go-vln is a vault symlink proxy`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	version = "local"
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.go-vln.yaml)")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Extra output")
	rootCmd.PersistentFlags().StringVarP(&backend, "backend", "b", "vault", "SymlinkDB backend")
	rootCmd.PersistentFlags().StringVarP(&filePath, "be-file-path", "f", "./db.json", "Path to file containing the KV Database in json.")
	rootCmd.PersistentFlags().StringVarP(&vaultAddr, "be-vault-addr", "a", "http://127.0.0.1:8200", "Vault ADDR to get secrets from.")
	rootCmd.PersistentFlags().StringVarP(&symlinkDBPath, "symlinkdb-path", "p", "secret/data/vln/symlinksdb", "Path to secret containing symlinks KV.")

	rootCmd.Flags().Bool("version", false, "Get Version")

	viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("backend", rootCmd.PersistentFlags().Lookup("backend"))
	viper.BindPFlag("be-file-path", rootCmd.PersistentFlags().Lookup("be-file-path"))
	viper.BindPFlag("be-vault-addr", rootCmd.PersistentFlags().Lookup("be-vault-addr"))
	viper.BindPFlag("symlinkdb-path", rootCmd.PersistentFlags().Lookup("symlinkdb-path"))

	viper.BindPFlags(rootCmd.Flags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-vln" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".go-vln")
	}
	viper.SetEnvPrefix("VLN")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

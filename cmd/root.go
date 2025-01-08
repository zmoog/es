package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zmoog/classeviva/adapters/feedback"
	"github.com/zmoog/es/cmd/docs"
	"github.com/zmoog/es/cmd/search"
	"github.com/zmoog/es/cmd/version"
)

var (
	cfgFile string
)

var rootCmd = cobra.Command{
	Use:   "es",
	Short: "Elasticsearch API via CLI",
	Long:  "Access the Elasticsearch API via CLI (currenly supports only a tiny subset of commands)",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if !viper.IsSet("api.endpoints") {
			must(cmd.MarkFlagRequired("api-endpoints"))
		}
		if !viper.IsSet("api.key") {
			must(cmd.MarkFlagRequired("api-key"))
		}

		if !viper.IsSet("client.max-retries") {
			must(cmd.MarkFlagRequired("max-retries"))
		}
		if !viper.IsSet("client.retry-on-status") {
			must(cmd.MarkFlagRequired("retry-on-status"))
		}
	},
	// Run: func(cmd *cobra.Command, args []string) {
	// 	feedback.Println("[" + viper.GetString("api.endpoints") + "]")
	// 	feedback.Println("[" + viper.GetString("api.key") + "]")
	// },
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		feedback.Error(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.es/config)")
	rootCmd.PersistentFlags().StringP("api-endpoints", "e", "", "Elasticsearch API endpoints")
	rootCmd.PersistentFlags().StringP("api-key", "k", "", "Elasticsearch API key")

	rootCmd.PersistentFlags().IntSliceP("client-retry-on-status", "r", []int{502, 503, 504, 429}, "Retry on status codes")
	rootCmd.PersistentFlags().IntP("client-max-retries", "m", 1, "Maximum number of retries")
	rootCmd.PersistentFlags().StringP("client-ca-cert-path", "c", "", "CA certificate path")

	rootCmd.AddCommand(docs.NewCommand())
	rootCmd.AddCommand(search.NewCommand())
	rootCmd.AddCommand(version.NewCommand())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		configDir := filepath.Join(home, ".es")

		// Search config in home directory with name ".es" (without extension).
		viper.AddConfigPath(configDir)
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	viper.SetEnvPrefix("ES")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	must(viper.BindPFlag("api.endpoints", rootCmd.PersistentFlags().Lookup("api-endpoints")))
	must(viper.BindPFlag("api.key", rootCmd.PersistentFlags().Lookup("api-key")))

	must(viper.BindPFlag("client.retry-on-status", rootCmd.PersistentFlags().Lookup("client-retry-on-status")))
	must(viper.BindPFlag("client.max-retries", rootCmd.PersistentFlags().Lookup("client-max-retries")))
	must(viper.BindPFlag("client.ca-cert-path", rootCmd.PersistentFlags().Lookup("client-ca-cert-path")))
}

func must(err error) {
	if err != nil {
		feedback.Error(err)
		os.Exit(1)
	}
}

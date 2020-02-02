package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/viper"
	server "github.com/import-yuefeng/BGPParser/daemon"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
	daemon 	 	bool

	rootCmd = &cobra.Command{
		Use:   "bgpParser",
		Short: "A parser for BGPData ",
		Long: `bgpParser is a CLI tool for Go that parse bdpData.`,
		Run: func(cmd *cobra.Command, args []string) {
			if daemon {
				server.Daemon()
			}
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.SuggestionsMinimumDistance = 1

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.bgpParser.json)")
	rootCmd.PersistentFlags().StringP("author", "a", "Yuefeng Zhu", "Copyright (c) 2019 Yuefeng Zhu")
	rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "MIT License")
	rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	rootCmd.PersistentFlags().BoolVar(&daemon, "daemon", false, "run BGPParser with daemon mode(default: client)")

	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	viper.SetDefault("author", "Yuefeng Zhu <zhuyuefeng0@gmail.com>")
	viper.SetDefault("license", "MIT License")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".cobra")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		log.Info("Using config file:", viper.ConfigFileUsed())
	}
}
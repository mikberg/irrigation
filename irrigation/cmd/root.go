package cmd

import (
	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Used for flags
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "irrigation",
		Short: "Irrigation is cool",
		Run: func(cmd *cobra.Command, args []string) {
			runCmd.Run(cmd, args)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.irrigation.yaml)")

	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(analogCmd)
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name '.irrigation'
		viper.AddConfigPath(home)
		viper.SetConfigName(".irrigation")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Info().Msgf("using config file: %s", viper.ConfigFileUsed())
	}
}

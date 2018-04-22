package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "medkit",
	Short:   "MEDKIT is a multi-environment dotfiles manager",
	Long:    `MEDKIT (Multi-Environment Dotfiles Kit) is the dotfile management solution for the obsessive compulsive.`,
	Version: "0.0.1-alpha",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// init initializes the command.
func init() {
	cobra.OnInitialize(initConfig)
	initFlags()
	setDefaults()
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".medkit" (without extension).
		viper.AddConfigPath(getHome())
		viper.SetConfigName(".medkit")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}

// initFlags initializes the flags.
func initFlags() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here, will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile,"config", "", "config file (default is $HOME/.medkit)")
}

// setDefaults sets the program defaults.
func setDefaults() {
	home := getHome()
	viper.SetDefault("HomeDirectory", home)
	viper.SetDefault("DotfilesDirectory", home+"/dotfiles")
	viper.SetDefault("Bundles", "")
	viper.SetDefault("BackupExtension", ".backup")
}

// getHome returns the home directory.
func getHome() string {
	home, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}

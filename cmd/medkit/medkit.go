package medkit

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var Version string
var BuildDate string
var CommitHash string

var cfgFile string

// medkitCmd represents the base command when called without any sub-commands.
var medkitCmd = &cobra.Command{
	Use:     "medkit",
	Short:   "MEDKIT is a multi-environment dotfiles manager",
	Long:    `MEDKIT (Multi-Environment Dotfiles Kit) is the dotfile management solution for the obsessive compulsive.`,
	Version: Version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the medkitCmd.
func Execute() {
	if err := medkitCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	medkitCmd.SetVersionTemplate(fmt.Sprintf("%s\n%s (%s)", Version, BuildDate, CommitHash))

	// Global flags
	medkitCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.medkit)")

	// Local flags
	medkitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	setDefaults()

	// Commands
	medkitCmd.AddCommand(installCmd)
	medkitCmd.AddCommand(showCmd)
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

	// If a config file is found, read it in
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file: ", viper.ConfigFileUsed())
	}
}

// setDefaults sets the program's default settings.
func setDefaults() {
	home := getHome()
	viper.SetDefault("HomeDirectory", home)
	viper.SetDefault("DotfilesDirectory", home+"/dotfiles")
	viper.SetDefault("Bundles", "")
	viper.SetDefault("BackupExtension", ".backup")
}

// getHome returns the current user's home directory.
func getHome() string {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return home
}

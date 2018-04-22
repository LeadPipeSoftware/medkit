package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// showConfigCmd represents the showConfig command
var showConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current MEDKIT configuration",
	Long:  `Display current MEDKIT configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MEDKIT's current configuration values:")
		fmt.Println(" homeDirectory: " + viper.GetString("homeDirectory"))
		fmt.Println(" dotfilesDirectory: " + viper.GetString("dotfilesDirectory"))
		fmt.Println(" bundles: " + viper.GetString("bundles"))
	},
}

func init() {
	showCmd.AddCommand(showConfigCmd)
}

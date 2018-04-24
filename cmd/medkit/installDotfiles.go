package medkit

import (
	"github.com/spf13/cobra"
	"github.com/LeadPipeSoftware/medkit/internal/dotfile"
)

var alwaysSkip = false
var alwaysOverwrite = false

var installDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "install .symlink files in your dotfiles directory",
	Long: `
This command looks for any file with a .symlink extension in your dotfiles
directory. When it finds a match, it will create a symbolic link from that
file to your home directory.`,
	Run: dotfile.InstallDotfiles,
}

func init() {
	installDotfilesCmd.Flags().BoolVarP(&alwaysSkip, "always-skip", "s", false, "always skip existing files")
	installDotfilesCmd.Flags().BoolVarP(&alwaysOverwrite, "always-overwrite", "o", false, "always overwrite existing files")
}


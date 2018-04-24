package medkit

import (
	"github.com/spf13/cobra"
	"github.com/LeadPipeSoftware/medkit/internal/dotfile"
)

// showDotfilesCmd represents the showDotfiles command
var showDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "display the dotfiles in your dotfiles directory",
	Long: `
Display all of the dotfiles contained in your dotfiles directory. This command
is recursive.`,
	Run: dotfile.ShowDotfiles,
}

func init() {
}

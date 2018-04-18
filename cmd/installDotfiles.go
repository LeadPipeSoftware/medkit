package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const bundlesDir = "bundles"

var alwaysSkip bool
var alwaysOverwrite bool

var installDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "Creates symbolic links based on the .symlink files in your dotfiles directory.",
	Long: `This command looks for any file with a .symlink extension in your dotfiles directory. When it finds a match,
it will create a symbolic link from that file to your home directory.
    
Run this command any time you want to make sure all your dotfiles have been installed.`,
	Run: func(cmd *cobra.Command, args []string) {

		if alwaysSkip && alwaysOverwrite {
			fmt.Println("The always-skip and always-overwrite flags cannot be used together.")
			os.Exit(1)
		}

		dotfilesDirectory := viper.GetString("DotfilesDirectory")

		fmt.Printf("Dotfile directory: %s\n", dotfilesDirectory)
		fmt.Printf("   Home directory: %s\n", viper.GetString("HomeDirectory"))
		fmt.Println()

		if err := filepath.Walk(dotfilesDirectory, visit); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	installDotfilesCmd.Flags().BoolVarP(&alwaysSkip, "always-skip", "s", false, "always skip existing files")
	installDotfilesCmd.Flags().BoolVarP(&alwaysSkip, "always-overwrite", "o", false, "always overwrite existing files")
	installCmd.AddCommand(installDotfilesCmd)
}

// The function called by the Walk method.
func visit(path string, f os.FileInfo, err error) error {
	home := viper.GetString("HomeDirectory")

	if f.IsDir() && f.Name() == bundlesDir {
		return filepath.SkipDir
	}

	if f.IsDir() {
		if err := symlinkFilesInDirectory(path, home); err != nil {
			return err
		}
	}

	return nil
}

// Creates a backup of and then removes a file.
func backupAndRemoveTarget(filename string) error {
	backupExtension := viper.GetString("BackupExtension")
	backupFile := filename + backupExtension

	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		fmt.Println("No existing backup")
	} else {
		if err := os.Remove(backupFile); err == nil {
			fmt.Println("Existing backup removed")
		} else {
			fmt.Printf("ERROR removing backup %s: %s\n", backupFile, err)
			return err
		}
	}

	if err := os.Rename(filename, backupFile); err == nil {
		fmt.Println("Backup created")
	} else {
		fmt.Printf("ERROR creating backup %s: %s\n", backupFile, err)
		return err
	}

	return nil
}

// Creates a symbolic link.
func createSymlink(source string, target string) {
	if err := os.Symlink(source, target); err == nil {
		fmt.Printf("Symlinked %s => %s\n", source, target)
	} else {
		fmt.Printf("ERROR symlinking %s: %s\n", target, err)
	}
}

// Returns true if the specified file does not exist.
func fileDoesNotExist(targetFile string) bool {
	_, err := os.Stat(targetFile)

	return os.IsNotExist(err)
}

// Builds and returns the symbolic link target name (the name without the .symlink extension).
func getSymlinkTargetName(fileName string) string {
	name := path.Base(fileName)
	re := regexp.MustCompile("\\.symlink")

	return re.ReplaceAllString(name, "")
}

// Creates symbolic links for each .symlink file in a directory.
func symlinkFilesInDirectory(path string, home string) error {
	matches, err := filepath.Glob(path + "/*.symlink")
	if err == nil {
		for _, match := range matches {
			targetFile := home + "/" + getSymlinkTargetName(match)
			if fileDoesNotExist(targetFile) {
				createSymlink(match, targetFile)
			} else {
				input := bufio.NewScanner(os.Stdin)

				if alwaysSkip {
					fmt.Printf("Skipping %s\n", targetFile)
				} else if alwaysOverwrite {
					fmt.Printf("Overwriting %s\n", targetFile)
					if err := backupAndRemoveTarget(targetFile); err == nil {
						createSymlink(match, targetFile)
					}
				} else {
					fmt.Printf("%s already exists. (S)kip or (O)verwrite?\n", targetFile)

				InputLoop:
					for input.Scan() {

						answer := strings.ToLower(input.Text())

						switch answer {
						case "o":
							fmt.Printf("Overwriting %s\n", targetFile)
							if err := backupAndRemoveTarget(targetFile); err == nil {
								createSymlink(match, targetFile)
							}
							break InputLoop
						case "s":
							fmt.Printf("Skipping %s\n", targetFile)
							break InputLoop
						default:
							fmt.Println("Invalid response. (S)kip or (O)verwrite?")
						}
					}
				}
			}
		}
	} else {
		return err
	}

	return nil
}

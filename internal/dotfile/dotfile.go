package dotfile

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const bundlesDir = "bundles"

var AlwaysSkip bool
var AlwaysOverwrite bool

// InstallDotfiles will install dotfiles.
func InstallDotfiles(cmd *cobra.Command, args []string) {
	if AlwaysSkip && AlwaysOverwrite {
		fmt.Print("Sorry, the always-skip and always-overwrite flags cannot be used together.")
		os.Exit(1)
	}

	dotfilesDirectory := viper.GetString("DotfilesDirectory")

	fmt.Printf("Dotfile directory: %s", dotfilesDirectory)
	fmt.Printf("\n   Home directory: %s", viper.GetString("HomeDirectory"))
	fmt.Println()

	if err := filepath.Walk(dotfilesDirectory, visit); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println()
}

// visit determines the action(s) taken during the filepath.Walk method.
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

// backupThenRemoveFile creates a backup of a file, but removes existing backups first.
func backupThenRemoveFile(filename string) error {
	backupExtension := viper.GetString("BackupExtension")
	backupFile := filename + backupExtension

	if fileExists(backupFile) {
		if err := os.Remove(backupFile); err != nil {
			fmt.Printf("\nERROR removing backup %s: %s", backupFile, err)
			return err
		}
	}

	if err := os.Rename(filename, backupFile); err != nil {
		fmt.Printf("\nERROR creating backup %s: %s", backupFile, err)
		return err
	}

	return nil
}

// createSymlink creates a symbolic link.
func createSymlink(source string, target string) {
	if err := os.Symlink(source, target); err == nil {
		fmt.Printf("\nSymlinked %s => %s", source, target)
	} else {
		fmt.Printf("\nERROR symlinking %s: %s", target, err)
	}
}

// fileDoesNotExist returns true if the specified file does not exist.
func fileDoesNotExist(targetFile string) bool {
	_, err := os.Stat(targetFile)

	return os.IsNotExist(err)
}

// fileExists returns true if the specified file exists.
func fileExists(targetFile string) bool {
	return !fileDoesNotExist(targetFile)
}

// getSymlinkTargetName builds and returns the symbolic link target name (the name without the .symlink extension).
func getSymlinkTargetName(fileName string) string {
	name := path.Base(fileName)
	re := regexp.MustCompile("\\.symlink")

	return re.ReplaceAllString(name, "")
}

// symlinkFilesInDirectory creates symbolic links for each .symlink file in a directory.
func symlinkFilesInDirectory(path string, home string) error {
	matches, err := filepath.Glob(path + "/*.symlink")

	if err == nil {
		for _, match := range matches {
			targetFile := home + "/" + getSymlinkTargetName(match)

			if fileDoesNotExist(targetFile) {
				createSymlink(match, targetFile)
			} else {
				input := bufio.NewScanner(os.Stdin)

				if AlwaysSkip {
					fmt.Printf("\nSkipping %s", targetFile)
				} else if AlwaysOverwrite {
					fmt.Printf("\nOverwriting %s", targetFile)
					if err := backupThenRemoveFile(targetFile); err == nil {
						createSymlink(match, targetFile)
					}
				} else {
					fmt.Printf("\n\n%s already exists. What do you want to do?\n[s]kip, [S]kip All, [o]verwrite, [O]verwrite All: ", targetFile)

				InputLoop:
					for input.Scan() {

						answer := input.Text()

						switch answer {
						case "O", "o":
							fmt.Printf("\nOkay, overwriting %s", targetFile)
							if err := backupThenRemoveFile(targetFile); err == nil {
								createSymlink(match, targetFile)
							}
							if answer == "O" {
								AlwaysOverwrite = true
							}
							break InputLoop
						case "S", "s":
							fmt.Printf("\nOkay, skipping %s", targetFile)
							if answer == "S" {
								AlwaysSkip = true
							}
							break InputLoop
						default:
							fmt.Printf("\n%s already exists. What do you want to do?\n[s]kip, [S]kip All, [o]verwrite, [O]verwrite All: ", targetFile)
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

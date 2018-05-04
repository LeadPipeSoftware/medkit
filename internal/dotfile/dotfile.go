package dotfile

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"unicode"
)

const bundlesDir = "bundles"

var AlwaysSkip bool
var AlwaysOverwrite bool
var ForceReinstall bool
var HomeDirectory string

// InstallDotfiles will install dotfiles.
func InstallDotfiles(dotfilesDirectory string, homeDirectory string, alwaysSkip bool, alwaysOverwrite bool, forceReinstall bool, backupExtension string) {
	AlwaysSkip = alwaysSkip
	AlwaysOverwrite = alwaysOverwrite
	ForceReinstall = forceReinstall

	if AlwaysSkip && AlwaysOverwrite {
		fmt.Print("Sorry, the always-skip and always-overwrite flags cannot be used together.")
		os.Exit(1)
	}

	fmt.Printf("Dotfile directory: %s", dotfilesDirectory)
	fmt.Printf("\n   Home directory: %s", homeDirectory)
	fmt.Println()

	if err := filepath.Walk(dotfilesDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && info.Name() == bundlesDir {
			return filepath.SkipDir
		}

		if info.IsDir() {
			if err := symlinkFilesInDirectory(path, homeDirectory, backupExtension); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		fmt.Printf("\n%s\n", err)
		os.Exit(1)
	}

	fmt.Println()
}

// ShowDotfiles will display all dotfiles (recursively) in the dotfiles directory.
func ShowDotfiles(dotfilesDirectory string) {
	if dotfiles, err := getAllDotfiles(dotfilesDirectory); err == nil {
		for _, dotfile := range dotfiles {
			fmt.Printf("\n%s", dotfile)
		}
	} else {
		fmt.Printf("\n%s", err)
	}

	fmt.Println()
}

// getAllDotfiles returns the paths of every dotfile in a directory (recursively).
func getAllDotfiles(rootpath string) ([]string, error) {
	list := make([]string, 0, 10)

	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) == ".symlink" {
			list = append(list, path)
		}

		return nil
	})

	return list, err
}

// backupThenRemoveFile creates a backup of a file, but removes existing backups first.
func backupThenRemoveFile(fileName string, backupExtension string) error {
	backupFile := fileName + backupExtension

	if fileExists(backupFile) {
		if err := os.Remove(backupFile); err != nil {
			fmt.Printf("\nERROR removing backup %s: %s", backupFile, err)
			return err
		}
	}

	if err := os.Rename(fileName, backupFile); err != nil {
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

// fileIsSymlink returns true if the specified file is a symlink.
func fileIsSymlink(fileName string) (bool, error) {
	// Gotta use os.Lstat because os.Stat would read the target and not the link itself
	fi, err := os.Lstat(fileName)

	if err != nil {
		return false, err
	}

	if fi.Mode()&os.ModeSymlink != 0 {
		_, err := os.Readlink(fileName)

		if err != nil {
			return false, err
		}

		return true, nil
	} else {
		return false, nil
	}
}

// getSymlinkTargetName builds and returns the symbolic link target name (the name without the .symlink extension).
func getSymlinkTargetName(fileName string) string {
	name := path.Base(fileName)
	re := regexp.MustCompile("(?i)\\.symlink")

	return re.ReplaceAllString(name, "")
}

// getCaseInsensitiveFilePath returns a Glob pattern that is case insensitive (depending on OS).
func getCaseInsensitiveFilePath(path string) string {
	if runtime.GOOS == "windows" {
		return path
	}

	p := ""

	for _, r := range path {
		if unicode.IsLetter(r) {
			p += fmt.Sprintf("[%c%c]", unicode.ToLower(r), unicode.ToUpper(r))
		} else {
			p += string(r)
		}
	}

	return p
}

// symlinkFilesInDirectory creates symbolic links for each .symlink file in a directory.
func symlinkFilesInDirectory(path string, home string, backupExtension string) error {
	caseInsensitiveFilePath := getCaseInsensitiveFilePath(path + "/*.symlink")

	matches, err := filepath.Glob(caseInsensitiveFilePath)

	if err == nil {
		for _, match := range matches {
			targetFile := home + "/" + getSymlinkTargetName(match)

			if fileDoesNotExist(targetFile) {
				createSymlink(match, targetFile)
			} else {
				if isLink, _ := fileIsSymlink(targetFile); isLink == true {
					fileLinksTo, _ := os.Readlink(targetFile)

					if fileLinksTo == match {
						if ForceReinstall == false {
							fmt.Printf("\n%s already installed", targetFile)
							continue
						}
					}
				}
				if AlwaysSkip {
					fmt.Printf("\nSkipping %s", targetFile)
				} else if AlwaysOverwrite {
					fmt.Printf("\nOverwriting %s", targetFile)
					if err := backupThenRemoveFile(targetFile, backupExtension); err == nil {
						createSymlink(match, targetFile)
					}
				} else {
					fmt.Printf("\n\n%s already exists. What do you want to do?\n[s]kip, [S]kip All, [o]verwrite, [O]verwrite All: ", targetFile)

					input := bufio.NewScanner(os.Stdin)

				InputLoop:
					for input.Scan() {

						answer := input.Text()

						switch answer {
						case "O", "o":
							fmt.Printf("\nOkay, overwriting %s", targetFile)
							if err := backupThenRemoveFile(targetFile, backupExtension); err == nil {
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

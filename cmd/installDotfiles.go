// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

const bundlesDir string = "bundles"

var installDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "THE command. Sets up all your stuff.",
	Long: `This is the command that does all the things. It symlinks your *.symlink files, 
    sources your *.source files, paths your *.path files, and installs your *.installer files.  
    
    Run this command any time you have made changes to your dotfiles repo.`,
	Run: func(cmd *cobra.Command, args []string) {
		dotfilesDirectory := viper.GetString("dotfilesDirectory")

		fmt.Printf("Dotfile directory: %s\n", dotfilesDirectory)
		fmt.Printf("   Home directory: %s\n", viper.GetString("homeDirectory"))
		fmt.Println()

		if err := filepath.Walk(dotfilesDirectory, visit); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	},
}

func init() {
	installCmd.AddCommand(installDotfilesCmd)
}

func visit(path string, f os.FileInfo, err error) error {
	home := viper.GetString("homeDirectory")

	if f.IsDir() && f.Name() == bundlesDir {
		return filepath.SkipDir
	}

	if f.IsDir() {
		matches, err := filepath.Glob(path + "/*.symlink")
		if err == nil {
			for _, match := range matches {
				targetFile := home + "/" + getSymlinkTargetName(match)
				if fileDoesNotExist(targetFile) {
					createSymlink(match, targetFile)
				} else {
					input := bufio.NewScanner(os.Stdin)

					fmt.Printf("%s already exists. (S)kip or (O)verwrite?\n", targetFile)

					InputLoop:
					for input.Scan() {

						answer := strings.ToLower(input.Text())

						switch answer {
						case "o":
							fmt.Printf("Overwriting %s\n", f.Name())
							if err := backupAndRemoveTarget(targetFile); err == nil {
								createSymlink(match, targetFile)
							}
							break InputLoop
						case "s":
							fmt.Printf("Skipping %s\n", f.Name())
							break InputLoop
						default:
							fmt.Println("Invalid response. (S)kip or (O)verwrite?")
						}
					}
				}
			}
		} else {
			return err
		}
	}

	return nil
}

func backupAndRemoveTarget(filename string) error {
	backupExtension := ".backup"
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

func createSymlink(source string, target string) {
	if err := os.Symlink(source, target); err == nil {
		fmt.Printf("Symlinked %s => %s\n", source, target)
	} else {
		fmt.Printf("ERROR symlinking %s: %s\n", target, err)
	}
}

// Attempt to get the FileInfo for the target and if there's an error then assume the file doesn't exist
func fileDoesNotExist(targetFile string) bool {
	_, err := os.Stat(targetFile)

	return os.IsNotExist(err)
}

func getSymlinkTargetName(fileName string) string {
	name := path.Base(fileName)
	re := regexp.MustCompile("\\.symlink")

	return re.ReplaceAllString(name, "")
}

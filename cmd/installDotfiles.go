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
	"fmt"
    "os"
    "path"
    "path/filepath"
    "regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const bundlesDir string = "bundles"

// installDotfilesCmd represents the installDotfiles command
var installDotfilesCmd = &cobra.Command{
	Use:   "dotfiles",
	Short: "THE command.  Sets up all your stuff.",
	Long: `This is the command that does all the things.  It symlinks your *.symlink files, 
    sources your *.source files, paths your *.path files, and installs your *.installer files.  
    
    Run this command any time you have made changes to your dotfiles repo`,
	Run: func(cmd *cobra.Command, args []string) {
        home := viper.GetString("dotfilesDirectory")
		fmt.Println("Will process all the dotfiles in " + home)

        err := filepath.Walk(home, visit)
        fmt.Printf("visit returned: %v\n", err)
	},
}

func init() {
	installCmd.AddCommand(installDotfilesCmd)
}

func visit(path string, f os.FileInfo, err error) error {
    if f.IsDir() && f.Name() == bundlesDir {
        return filepath.SkipDir
    }
    if f.IsDir() {
        matches, err := filepath.Glob(path + "/*.symlink")
        if err == nil {
            for _, match := range matches {
                home := viper.GetString("homeDirectory")
                targetFile := home + "/" + getSymlinkTargetName(match)
                if shouldLink(targetFile) {
                    fmt.Println("Will symlink " + match + " => " + targetFile)
                }
            }
        } else {
            return err
        }
    }

    return nil
}

func getSymlinkTargetName(fileName string) string {
    name := path.Base(fileName)
    re := regexp.MustCompile("\\.symlink")

    return re.ReplaceAllString(name, "")
}

func shouldLink(targetFile string) bool {
    _, err := os.Stat(targetFile)

    return os.IsNotExist(err)
}

/*
Copyright © 2022 Fimreal

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"os"

	"github.com/fimreal/gitops-tools/config"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	Verbose *bool
	version *bool

	rootCmd = &cobra.Command{
		Use:   "gitops-tools",
		Short: config.AppName + "is a gitops tools",
		Long: config.AppName + `

Source code & Details:
Github: https://github.com/fimreal/gitops-tools`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		Run: func(cmd *cobra.Command, args []string) {
			if *version {
				ezap.Println(config.AppVersion)
			}
			if *Verbose {
				ezap.SetLevel("debug")
			}
			ezap.SetLogTime("")
			ezap.Println(config.AppName, config.AppVersion)
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitops-tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	Verbose = rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable debug log")
	version = rootCmd.Flags().BoolP("version", "", false, "Show this application version")
}

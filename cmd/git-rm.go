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
	"github.com/fimreal/gitops-tools/pkg/gitlab"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "删除的 git 仓库文件",
	Long: `删除 git 仓库文件. For example:

Alias command:

    alias gittool='` + Gittool + `'

Delete remote file:

    ` + Gittool + ` rm file1 file2 file3

`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化部分
		ezap.SetLogTime("")
		if *Verbose {
			ezap.SetLevel("debug")
			if len(args) == 0 {
				ezap.Error("未找到需要删除的文件名，请使用 -h 查看命令帮助")
			}
			for i, v := range args {
				ezap.Debugf("接收到参数: %d:%s", i, v)
			}
		}

		// set git provider, eg. gitlab github gitea
		if *GitProvider != "gitlab" {
			ezap.Fatal("Git provider [%s] is not support now.", *GitProvider)
		}
		gitlab.SetProvider(*GitProvider)

		// 配置连接
		gc, err := gitlab.ParseGcc(*GitClientConfig)
		if err != nil {
			ezap.Fatal(err)
		}

		for _, file := range args {
			ezap.Infof("删除文件: %v", file)
			err = gc.FileDelete(file, *GitCommitMessage)
			if err != nil {
				ezap.Error(err)
			}
		}
	},
}

func init() {
	gitCmd.AddCommand(rmCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rmCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// files = rmCmd.Flags().StringArrayP("files", "f", []string{}, "指定需要删除的文件，可多次使用")
}

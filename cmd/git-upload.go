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

// uploadCmd represents the upload command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "从 git 仓库复制文件，或者从本地复制文件到 git 仓库中",
	Long: `从 git 仓库复制文件，或者从本地复制文件到 git 仓库中, For example:

Alias command:

	alias gittool='` + Gittool + `'

upload local file[file1.txt] => remote file[file.txt]:

	` + Gittool + ` upload file1.txt file.txt

`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化部分
		ezap.SetLogTime("")
		if *Verbose {
			ezap.SetLevel("debug")
			if len(args) != 2 {
				ezap.Error("参数不正确，请分别输入本地文件名和远程文件名，也可以使用 -h 查看命令帮助")
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

		local := args[0]
		remote := args[1]
		ezap.Debugf("上传文件: %s => %s", local, remote)
		err = gc.Upload(local, remote, *GitCommitMessage)
		if err != nil {
			ezap.Error(err)
		}
		ezap.Info("上传成功")
	},
}

func init() {
	gitCmd.AddCommand(uploadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verbose = uploadCmd.Flags().BoolP("verbose", "v", false, "打开 debug 日志")
}

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

var (
	raw            *bool
	outputDocument *string
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get remote file from git project.",
	Long: `Get remote file from git project. For example:

Alias command:

	alias gittool='` + Gittool + `'

Get file[README.md] info:

    ` + Gittool + `get README.md

Download file[README.md] to file.txt:

    ` + Gittool + `get README.md -o file.txt

Get file[README.md] raw content:

    ` + Gittool + `get README.md -r

`,
	Run: func(cmd *cobra.Command, args []string) {
		// 初始化部分
		ezap.SetLogTime("")
		if *Verbose {
			ezap.SetLevel("debug")
			if len(args) == 0 {
				ezap.Error("未找到需要获取的文件名，请使用 -h 查看命令帮助")
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

		if *outputDocument != "" {
			err = gc.Download(args[0], *outputDocument)
			if err != nil {
				ezap.Fatal(err)
			}
			ezap.Info(*outputDocument, " 下载完成!")
			return
		}

		// 判断是否只获取 raw 内容
		if *raw {
			for _, file := range args {
				c, err := gc.GetFileRaw(file)
				if err != nil {
					ezap.Error(err)
				} else {
					ezap.Println(string(c))
				}
			}
		} else {
			for _, file := range args {
				c, err := gc.GetFileInfo(file)
				if err != nil {
					ezap.Error(err)
				} else {
					ezap.Println(c)
				}
			}
		}

	},
}

func init() {
	gitCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	raw = getCmd.Flags().BoolP("raw", "r", false, "Only request raw content.")
	outputDocument = getCmd.Flags().StringP("output-document", "o", "", "将文件另存为, 参数添加下载保存的文件名。只接受接收第一个参数进行下载")
}

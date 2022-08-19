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
	"github.com/fimreal/gitops-tools/config"
	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

var (
	GitClientConfig  *string
	GitProvider      *string
	GitCommitMessage *string

	Gittool = config.AppName + " git -c \"http://lxm:token@www.gitlab.cn/lxm/project-1?branch=master&mail=lxm@mail.com\""
)

// gitCmd represents the git command
var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "git 命令可实现远程仓库控制相关操作",
	Long: `git 命令可实现远程仓库控制相关操作。
`,
	Run: func(cmd *cobra.Command, args []string) {
		ezap.Warn("添加 -h 参数查看使用命令帮助")
	},
}

func init() {
	rootCmd.AddCommand(gitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gitCmd.PersistentFlags().String("foo", "", "A help for foo")
	GitProvider = gitCmd.PersistentFlags().StringP("provider", "p", "gitlab", "指定 git provider 类型, (目前只支持)默认使用 gitlab")
	GitClientConfig = gitCmd.PersistentFlags().StringP("client-config", "c", "", "git 连接配置串，具体到仓库项目名和分支，例如 "+"http://lxm:token@www.gitlab.cn/lxm/project-1?branch=master&mail=gitops@mail.com&useragent=gohttp "+"参数可选，例子中为默认值")
	GitCommitMessage = gitCmd.PersistentFlags().StringP("commit-message", "m", "commit by gitops-tools, please use --commit-message set this message", "自定义 git commit 备注信息")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

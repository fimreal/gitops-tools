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
	"fmt"

	"github.com/fimreal/goutils/ezap"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var (
	// 更新镜像的容器名字
	container *string
	file      *string

	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "更新 kubernetes yaml 配置中镜像",
		Long: `
	用法说明：
	
	gitops-tools set image -f demo.yaml -c [container name] <image name>
	
	`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("image called")

			if len(args) != 1 {
				ezap.Fatal("请指定镜像名字")
			}

			if *container == "" {
				ezap.Warn("未指定容器名字，修改配置文件中第一个镜像")
				*container = "containerName"
			}

			// do
			ezap.Infof("成功将文件[%s]中容器[%s]镜像修改为: %s", *file, *container, args[0])

		},
	}
)

func init() {
	setCmd.AddCommand(imageCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// imageCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// imageCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	file = imageCmd.Flags().StringP("file", "f", "", "指定对其操作的文件")
	imageCmd.MarkFlagRequired("file")
	container = imageCmd.Flags().StringP("container", "c", "", "更新使用的镜像, 默认更新文件中第一个镜像")
}

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
	"github.com/fimreal/gitops-tools/pkg/yaml"

	"github.com/fimreal/goutils/ezap"
	mfile "github.com/fimreal/goutils/file"
	"github.com/spf13/cobra"
)

// imageCmd represents the image command
var (
	// 更新镜像的容器名字
	containerName *string
	file          *string

	imageCmd = &cobra.Command{
		Use:   "image",
		Short: "更新 kubernetes yaml 配置中镜像",
		Long: `
用法举例：

  gitops-tools set image -f demo.yaml -c [container name] <image name>
	
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

			if len(args) != 1 {
				ezap.Fatal("请指定更新使用的镜像名字")
			}

			// do
			yamlByte, err := yaml.SplitYamlFile(*file)
			if err != nil {
				ezap.Fatal(err)
			}

			var newFileByte [][]byte
			for _, y := range yamlByte {
				err := y.UpdateImage(args[0], *containerName)
				if err != nil {
					ezap.Fatal(err)
				}
				newFileByte = append(newFileByte, y.ByteData, []byte("\n---\n"))
			}

			// 清空文件
			mfile.WriteToFile(*file, nil)
			// 写入文件
			for _, frag := range newFileByte {
				mfile.AppendToFile(*file, frag)
			}

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
	containerName = imageCmd.Flags().StringP("container", "c", "", "更新使用的镜像, 默认更新文件中第一个镜像")
}

package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/jumpserver/replay_uploader/cmd/common"
	"github.com/jumpserver/replay_uploader/cmd/scan"
	"github.com/jumpserver/replay_uploader/util"
)

/*
	包含录像文件的目录
	前置条件:
		文件名包含 Session ID
	相关逻辑:
		存储需要日期
			如果会话文件父目录是日期格式，date则选择该日期，否则日期选择文件修改的时间
*/

var scanCommand = &cobra.Command{
	Use:   "scan",
	Short: "scan dir to upload replay files",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err            error
			plainAccessKey string
		)
		if accessKeyFile != "" {
			result, err := ioutil.ReadFile(accessKeyFile)
			if err != nil {
				msg := fmt.Sprintf("读取 access key 文件失败 %s", accessKeyFile)
				common.ReturnErrorMsg(msg, err)
				return
			}
			plainAccessKey = string(bytes.TrimSpace(result))
		}
		if accessKey != "" {
			result, err := util.DecodeBase64String(accessKey)
			if err != nil {
				msg := fmt.Sprintf("无法解析 base64 字符 %s", accessKey)
				common.ReturnErrorMsg(msg, err)
				return
			}
			plainAccessKey = result
		}

		jmsService, err := common.NewJmsAuthService(coreHost, plainAccessKey)
		if err != nil {
			msg := fmt.Sprintf("Core URL或认证信息失败: %s %s", coreHost, accessKey)
			common.ReturnErrorMsg(msg, err)
			return
		}
		terminalConfig, err := jmsService.GetTerminalConfig()
		if err != nil {
			msg := fmt.Sprintf("与JMS Core %s 获取配置失败", coreHost)
			common.ReturnErrorMsg(msg, err)
			return
		}
		scan.Execute(jmsService, &terminalConfig, baseDir, forceDelete)
	},
}

func init() {
	rootCmd.AddCommand(scanCommand)
	scanCommand.Flags().StringVar(&baseDir, "baseDir", "", "需要扫描的录像文件目录")
}

var (
	baseDir string
)

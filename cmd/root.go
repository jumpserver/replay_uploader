package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/jumpserver-dev/sdk-go/model"
	"github.com/jumpserver/replay_uploader/cmd/common"
	"github.com/jumpserver/replay_uploader/cmd/upload"
	"github.com/jumpserver/replay_uploader/util"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "replay_uploader",
	Short: "Jumpserver replay upload tool",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			err            error
			plainAccessKey string
		)
		if accessKeyFile != "" {
			result, err := os.ReadFile(accessKeyFile)
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

		if replayPath == "" {
			err := fmt.Errorf("未发现录像文件: %s", replayPath)
			common.ReturnErrorMsg("未发现录像文件", err)
			return
		}
		var (
			ok bool
		)
		if sid == "" {
			if sid, ok = common.ParseSessionID(replayPath); !ok {
				msg := fmt.Sprintf("无合法的会话ID: %s", replayPath)
				common.ReturnErrorMsg(msg, err)
				return
			}
		}
		if !util.IsUUID(sid) {
			msg := fmt.Sprintf("不是合法的会话ID %s", sid)
			err := fmt.Errorf("不是合法的会话ID %s", sid)
			common.ReturnErrorMsg(msg, err)
			return
		}
		if targetDate == "" {
			if targetDate, ok = common.ParseDateFromPath(replayPath); !ok {
				targetDate = util.CurrentDate()
			}
		}
		var replayVersion = model.Version2
		if version, ok := common.ParseReplayVersion(replayPath); ok {
			replayVersion = version
		}

		replayFile := common.ReplayFile{
			ID:          sid,
			TargetDate:  targetDate,
			AbsFilePath: replayPath,
			Version:     replayVersion,
		}

		if err := upload.Execute(jmsService, &terminalConfig,
			&replayFile); err != nil {
			common.ReturnErrorMsg("上传文件失败", err)
			return
		}
		if forceDelete {
			_ = os.Remove(replayPath)
		}
		msg := fmt.Sprintf("会话 %s 录像文件上传成功 %s", sid, replayPath)
		common.ReturnSuccessMsg(msg)
	},
}

var (
	coreHost      string
	accessKey     string
	accessKeyFile string
	forceDelete   bool
)

var (
	replayPath string
	targetDate string
	sid        string
)

func init() {
	rootCmd.PersistentFlags().StringVar(&coreHost, "url", "http://127.0.0.1:8080", "JumpServer URL，JMS Core的地址")
	rootCmd.PersistentFlags().StringVar(&accessKey, "key", "", "Access key，组件使用的认证Key")
	rootCmd.PersistentFlags().StringVar(&accessKeyFile, "keyfile", "", "Key file，存储 access key 的文件路径")
	rootCmd.PersistentFlags().BoolVar(&forceDelete, "remove", false, "成功上传后，是否删除原 file 文件")

	rootCmd.Flags().StringVar(&replayPath, "file", "", "Replay file path，录像文件路径")
	rootCmd.Flags().StringVar(&sid, "sid", "", "Session ID, 会话ID")
	rootCmd.Flags().StringVar(&targetDate, "date", "", "Target date, 默认当前日期，格式 2021-01-20")
}

/*
	CoreHost 		指 JumpServer 的地址 通常是 http://127.0.0.1:8080
	accessKey 		指 经过 Base64 的 access_key 格式 keyId:secretId
	accessKeyFile   指 保存了access key 的文件，内容格式 keyId:secretId
	replayPath		指 录像文件，未经gzip压缩的录像文件
	targetDate		指 会话生成日期，上传 存储时需要
	sid				指 会话的ID
	forceDelete		指 成功上传之后，是否删除录像 file 源文件
*/

package upload

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/jumpserver/replay_uploader/storage"
	"github.com/jumpserver/replay_uploader/util"
)

func Execute() {
	flag.Parse()
	if infoFlag {
		fmt.Printf("Version:             %s\n", Version)
		fmt.Printf("Git Commit Hash:     %s\n", GitHash)
		fmt.Printf("UTC Build Time :     %s\n", BuildStamp)
		fmt.Printf("Go Version:          %s\n", GoVersion)
		return
	}
	if targetDate == "" {
		targetDate = util.CurrentDate()
	}
	var err error
	if sid == "" {
		if sid, err = util.ParseSessionID(replayPath); err != nil {
			ReturnErrorMsg("无合法的会话ID", err)
			return
		}
	}
	if !util.IsUUID(sid) {
		msg := fmt.Sprintf("不是合法的会话ID %s", sid)
		err := fmt.Errorf("不是合法的会话ID %s", sid)
		ReturnErrorMsg(msg, err)
		return
	}
	if replayPath == "" {
		err := fmt.Errorf("未发现录像文件: %s", replayPath)
		ReturnErrorMsg("未发现录像文件", err)
		return
	}

	plainAccessKey := ""
	if accessKeyFile != "" {
		result, err := ioutil.ReadFile(accessKeyFile)
		if err != nil {
			msg := fmt.Sprintf("读取 access key 文件失败 %s", accessKeyFile)
			ReturnErrorMsg(msg, err)
			return
		}
		plainAccessKey = string(bytes.TrimSpace(result))
	}
	if accessKey != "" {
		result, err := util.DecodeBase64String(accessKey)
		if err != nil {
			msg := fmt.Sprintf("无法解析 base64 字符 %s", accessKey)
			ReturnErrorMsg(msg, err)
			return
		}
		plainAccessKey = result
	}
	jmsService, err := NewJmsAuthService(coreHost, plainAccessKey)
	if err != nil {
		msg := fmt.Sprintf("Core URL或认证信息失败: %s %s", coreHost, accessKey)
		ReturnErrorMsg(msg, err)
		return
	}
	terminalConfig, err := jmsService.GetTerminalConfig()
	if err != nil {
		msg := fmt.Sprintf("与JMS Core %s 获取配置失败", coreHost)
		ReturnErrorMsg(msg, err)
		return
	}
	dirPath := filepath.Dir(replayPath)
	sidReplayPath := filepath.Join(dirPath, sid+SuffixReplayFileName)

	if !util.IsGzipFile(replayPath) {
		if err = util.CompressToGzipFile(replayPath, sidReplayPath); err != nil {
			msg := fmt.Sprintf("压缩录像文件失败 %s", replayPath)
			ReturnErrorMsg(msg, err)
			return
		}
		defer os.Remove(sidReplayPath)
	} else {
		if replayPath != sidReplayPath {
			if err = util.CopyFile(replayPath, sidReplayPath); err != nil {
				msg := fmt.Sprintf("录像文件重命名失败 %s", replayPath)
				ReturnErrorMsg(msg, err)
				return
			}
			defer os.Remove(sidReplayPath)
		}
	}

	replayStorage := storage.NewReplayStorage(&terminalConfig)
	if replayStorage == nil {
		err = jmsService.Upload(sid, sidReplayPath)
	} else {
		target := strings.Join([]string{targetDate, sid + SuffixReplayFileName}, "/")
		err = replayStorage.Upload(sidReplayPath, target)
	}
	if err != nil {
		msg := fmt.Sprintf("上传文件失败 %s", sidReplayPath)
		ReturnErrorMsg(msg, err)
		return
	}
	if forceDelete {
		_ = os.Remove(replayPath)
	}
	err = jmsService.FinishReply(sid)
	if err != nil {
		ReturnErrorMsg("通知Core录像文件上传完成失败", err)
		return
	}
	msg := fmt.Sprintf("会话 %s 录像文件上传成功 %s", sid, replayPath)
	ReturnSuccessMsg(msg)
}

const SuffixReplayFileName = ".replay.gz"

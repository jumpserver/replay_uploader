package upload

import (
	"flag"
	"fmt"
	"jms-upload/pkg/model"

	"jms-upload/pkg/storage"
)

func Execute() {
	flag.Parse()
	if targetDate == "" {
		targetDate = model.GetCurrentDate()
	}
	sid, err := model.ParseSessionID(replayPath)
	if err != nil {
		msg := fmt.Sprintf("不是合法的录像文件格式 %s", replayPath)
		ReturnErrorMsg(msg, err)
		return
	}
	jmsService, err := NewJmsAuthService(coreHost, accessKey)
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

	replayStorage := storage.NewReplayStorage(&terminalConfig)
	if replayStorage == nil {
		err = jmsService.Upload(sid, replayPath)
	} else {
		err = replayStorage.Upload(replayPath, targetDate)
	}
	if err != nil {
		msg := fmt.Sprintf("上传文件失败 %s", replayPath)
		ReturnErrorMsg(msg, err)
		return
	}
	err = jmsService.FinishReply(sid)
	if err != nil {
		ReturnErrorMsg("通知Core录像文件上传完成失败", err)
		return
	}
	msg := fmt.Sprintf("上传成功 %s", replayPath)
	ReturnSuccessMsg(msg)
}

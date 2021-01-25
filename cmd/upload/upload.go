package upload

import (
	"flag"
	"fmt"

	"jms-upload/cmd/common"
	"jms-upload/pkg/storage"
)

func Execute() {
	flag.Parse()
	if targetDate == "" {
		targetDate = common.GetCurrentDate()
	}
	sid, err := common.ParseSessionID(replayPath)
	if err != nil {
		msg := fmt.Sprintf("不是合法的录像文件格式 %s", replayPath)
		common.ReturnErrorMsg(msg, err)
		return
	}
	jmsService, err := common.NewJmsAuthService(coreHost, accessKey)
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

	replayStorage := storage.NewReplayStorage(&terminalConfig)
	if replayStorage == nil {
		err = jmsService.Upload(sid, replayPath)
	} else {
		err = replayStorage.Upload(replayPath, targetDate)
	}
	if err != nil {
		msg := fmt.Sprintf("上传文件失败 %s", replayPath)
		common.ReturnErrorMsg(msg, err)
		return
	}
	err = jmsService.FinishReply(sid)
	if err != nil {
		common.ReturnErrorMsg("通知Core录像文件上传完成失败", err)
		return
	}
	msg := fmt.Sprintf("上传成功 %s", replayPath)
	common.ReturnSuccessMsg(msg)
}

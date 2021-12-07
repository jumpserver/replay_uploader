package upload

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jumpserver/replay_uploader/cmd/common"
	"github.com/jumpserver/replay_uploader/jms-sdk-go/model"
	"github.com/jumpserver/replay_uploader/jms-sdk-go/service"
	"github.com/jumpserver/replay_uploader/storage"
	"github.com/jumpserver/replay_uploader/util"
)

func Execute(jmsService *service.JMService, conf *model.TerminalConfig, replay *common.ReplayFile, successCallback func()) error {
	replayAbsGzPath := replay.AbsFilePath
	if !util.IsGzipFile(replayAbsGzPath) {
		dirPath := filepath.Dir(replay.AbsFilePath)
		replayAbsGzPath = filepath.Join(dirPath, replay.GetGzFilename())
		if err := util.CompressToGzipFile(replay.AbsFilePath, replayAbsGzPath); err != nil {
			return fmt.Errorf("压缩录像文件失败 %s: %s", replay.AbsFilePath, err)
		}
		defer os.Remove(replayAbsGzPath)
	}

	if replayStorage := storage.NewReplayStorage(conf); replayStorage != nil {
		if err := replayStorage.Upload(replayAbsGzPath, replay.TargetPath()); err != nil {
			return fmt.Errorf("上传文件失败 %s", err)
		}
	} else {
		if err := jmsService.UploadReplay(replay.ID, replayAbsGzPath, replay.Version); err != nil {
			return fmt.Errorf("上传文件失败 %s", err)
		}
	}
	if successCallback != nil {
		successCallback()
	}
	if err := jmsService.FinishReply(replay.ID); err != nil {
		return fmt.Errorf("通知Core录像文件上传完成失败: %s", err)
	}
	return nil
}

package scan

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/jumpserver/replay_uploader/cmd/common"
	"github.com/jumpserver/replay_uploader/cmd/upload"
	"github.com/jumpserver/replay_uploader/jms-sdk-go/model"
	"github.com/jumpserver/replay_uploader/jms-sdk-go/service"
)

func Execute(jmsService *service.JMService, conf *model.TerminalConfig, rootDir string, forceDelete bool) {
	if !haveDir(rootDir) {
		res := Response{
			Code: CodeErr,
			Err:  fmt.Sprintf("无效的根目录 %s", rootDir),
		}
		ReturnResp(os.Stderr, res)
		return
	}

	allReplays := make([]common.ReplayFile, 0, 20)
	_ = filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil && info.IsDir() {
			return nil
		}
		var (
			sid        string
			targetDate string
			version    model.ReplayVersion
			ok         bool
		)

		if sid, ok = common.ParseSessionID(path); !ok {
			return nil
		}
		if targetDate, ok = common.ParseDateFromPath(path); !ok {
			targetDate = info.ModTime().Format("2006-01-02")
		}
		if version, ok = common.ParseReplayVersion(info.Name()); !ok {
			version = model.Version2
		}

		replayFile := common.ReplayFile{
			ID:          sid,
			TargetDate:  targetDate,
			AbsFilePath: path,
			Version:     version,
		}
		allReplays = append(allReplays, replayFile)
		return nil
	})

	if len(allReplays) == 0 {
		res := Response{
			Code: CodeOK,
		}
		ReturnResp(os.Stdout, res)
		return
	}
	// DISABLED_DELAY_UPLOAD

	if os.Getenv("DISABLED_DELAY_UPLOAD") != "1" {
		time.Sleep(10 * time.Minute)
	}

	var (
		successFiles []string
		failureFiles []string
		FailureErrs  []string
	)

	code := CodeOK

	for i := range allReplays {
		item := &allReplays[i]
		if err := upload.Execute(jmsService, conf, item); err != nil {
			code = CodeErr
			FailureErrs = append(FailureErrs, err.Error())
			failureFiles = append(failureFiles, item.AbsFilePath)
			continue
		}
		successFiles = append(successFiles, item.AbsFilePath)
		if forceDelete {
			_ = os.Remove(item.AbsFilePath)
		}
	}
	res := Response{
		Code:         code,
		SuccessFiles: successFiles,
		FailureFiles: failureFiles,
		FailureErrs:  FailureErrs,
	}
	switch code {
	case CodeOK:
		ReturnResp(os.Stdout, res)
	case CodeErr:
		ReturnResp(os.Stderr, res)
	default:
		ReturnResp(os.Stdout, res)
	}

}

const (
	CodeOK  = 200
	CodeErr = 400
)

type Response struct {
	Code         int      `json:"code"`
	Err          string   `json:"error"`
	SuccessFiles []string `json:"success_files,omitempty"`
	FailureFiles []string `json:"failure_files,omitempty"`
	FailureErrs  []string `json:"failure_errs,omitempty"`
}

func (r Response) String() string {
	result, _ := json.Marshal(r)
	return string(result)
}

func ReturnResp(w io.Writer, res Response) {
	_, _ = fmt.Fprint(w, res)
}

func haveDir(file string) bool {
	fi, err := os.Stat(file)
	return err == nil && fi.IsDir()
}

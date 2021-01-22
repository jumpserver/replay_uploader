package common

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"

	"jmsupload/pkg/httplib"
	"jmsupload/pkg/service"
)

func NewJmsAuthService(coreHost, accessKey string) (*service.JMService, error) {
	sigAuth, err := decodeAccessKey(accessKey)
	if err != nil {
		return nil, err
	}
	return service.NewAuthJMService(
		service.JMSCoreHost(coreHost),
		service.JMSTimeOut(time.Minute),
		service.JMSAuthSign(sigAuth))
}

func decodeAccessKey(p string) (*httplib.SigAuth, error) {
	result, err := base64.StdEncoding.DecodeString(p)
	if err != nil {
		return nil, err
	}
	accessKeyResult := bytes.Split(result, []byte(":"))
	if len(accessKeyResult) != 2 {
		return nil, fmt.Errorf("unknow accesskey %s", result)
	}

	return &httplib.SigAuth{
		KeyID:    string(bytes.TrimSpace(accessKeyResult[0])),
		SecretID: string(bytes.TrimSpace(accessKeyResult[1])),
	}, nil
}

// 录像的文件必须是固定格式 sessionID.replay.gz; eg: 90b11402-39df-4ba9-b14a-2344b8585888.replay.gz
func ParseSessionID(replayFilePath string) (string, error) {
	fileName := filepath.Base(replayFilePath)
	sid := strings.TrimRight(fileName, replayFileNameSuffix)
	if _, err := uuid.FromString(sid); err != nil {
		return "", err
	}
	return sid, nil
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

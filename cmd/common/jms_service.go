package common

import (
	"fmt"
	"strings"
	"time"

	"github.com/jumpserver/replay_uploader/jms-sdk-go/service"
)

func NewJmsAuthService(coreHost, accessKey string) (*service.JMService, error) {
	keyID, secretID, err := decodeAccessKey(accessKey)
	if err != nil {
		return nil, err
	}
	return service.NewAuthJMService(
		service.JMSCoreHost(coreHost),
		service.JMSTimeOut(time.Minute),
		service.JMSAccessKey(keyID, secretID))
}

func decodeAccessKey(p string) (keyID, secretID string, err error) {
	accessKeyResult := strings.Split(p, ":")
	if len(accessKeyResult) != 2 {
		return "", "", fmt.Errorf("unknow accesskey %s", p)
	}
	keyID = strings.TrimSpace(accessKeyResult[0])
	secretID = strings.TrimSpace(accessKeyResult[1])
	return
}

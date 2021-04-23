package upload

import (
	"fmt"
	"strings"
	"time"

	"github.com/jumpserver/replay_uploader/pkg/httplib"
	"github.com/jumpserver/replay_uploader/pkg/service"
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
	accessKeyResult := strings.Split(p, ":")
	if len(accessKeyResult) != 2 {
		return nil, fmt.Errorf("unknow accesskey %s", p)
	}
	return &httplib.SigAuth{
		KeyID:    strings.TrimSpace(accessKeyResult[0]),
		SecretID: strings.TrimSpace(accessKeyResult[1]),
	}, nil
}

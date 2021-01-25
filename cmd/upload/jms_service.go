package upload

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"jms-upload/pkg/httplib"
	"jms-upload/pkg/service"
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


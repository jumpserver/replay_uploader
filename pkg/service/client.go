package service

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/jumpserver/replay_uploader/pkg/httplib"
	"github.com/jumpserver/replay_uploader/pkg/model"
)

var AccessKeyUnauthorized = errors.New("access key unauthorized")

const (
	orgHeaderKey = "X-JMS-ORG"
	orgHeaderValue= "ROOT"
)

func NewAuthJMService(opts ...Option) (*JMService, error) {
	defaultOption := option{
		CoreHost: "http://127.0.0.1:8080",
		TimeOut:  time.Minute,
	}
	for _, setter := range opts {
		setter(&defaultOption)
	}
	httpClient, err := httplib.NewClient(defaultOption.CoreHost, defaultOption.TimeOut)
	if err != nil {
		return nil, err
	}
	if defaultOption.sign != nil {
		httpClient.SetAuthSign(defaultOption.sign)
	}
	httpClient.SetHeader(orgHeaderKey, orgHeaderValue)
	return &JMService{authClient: httpClient}, nil
}

type JMService struct {
	authClient *httplib.Client
}

func (s *JMService) GetUserDetail(userID string) (user *model.User) {
	url := fmt.Sprintf(UserDetailURL, userID)
	_, err := s.authClient.Get(url, &user)
	if err != nil {
		log.Println(err)
	}
	return
}

func (s *JMService) GetProfile() (user *model.User, err error) {
	var res *http.Response
	res, err = s.authClient.Get(UserProfileURL, &user)
	if err != nil {
		log.Println(err)
	}
	if res != nil && res.StatusCode == http.StatusUnauthorized {
		return user, AccessKeyUnauthorized
	}
	return user, err
}

func (s *JMService) GetTerminalConfig() (model.Config, error) {
	var conf model.Config
	_, err := s.authClient.Get(TerminalConfigURL, &conf)
	return conf, err
}

func (s *JMService) Upload(sessionID, gZipFile string) error {
	var res map[string]interface{}
	Url := fmt.Sprintf(SessionReplayURL, sessionID)
	return s.authClient.UploadFile(Url, gZipFile, &res)
}

func (s *JMService) FinishReply(sid string) error {
	var res map[string]interface{}
	data := map[string]bool{"has_replay": true}
	Url := fmt.Sprintf(SessionDetailURL, sid)
	_, err := s.authClient.Patch(Url, data, &res)
	return err
}

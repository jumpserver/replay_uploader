package service

import (
	"testing"
	"time"

	"github.com/jumpserver/replay_uploader/pkg/httplib"
)

func setupJMS() *JMService {
	CoreHost := "http://127.0.0.1:8080"
	authClient, _ := httplib.NewClient(CoreHost, 30*time.Second)
	sigAuth := httplib.SigAuth{
		KeyID:    "eaea44ec-cf0c-4bb4-98fa-ae2d2d140daf",
		SecretID: "43bb1290-a29b-4896-9916-a1756bd2fca2",
	}
	authClient.SetAuthSign(&sigAuth)
	return &JMService{
		authClient: authClient,
	}
}

func TestJMService_GetTerminalConfig(t *testing.T) {
	jms := setupJMS()
	result, err := jms.GetTerminalConfig()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

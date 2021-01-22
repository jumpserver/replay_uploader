package httplib

import (
	"log"
	"net/http"

	"gopkg.in/twindagger/httpsig.v1"
)

var (
	_ AuthSign = (*SigAuth)(nil)

	_ AuthSign = (*BasicAuth)(nil)
)

const (
	signHeaderRequestTarget = "(request-target)"
	signHeaderDate          = "date"
	signAlgorithm           = "hmac-sha256"
)

type SigAuth struct {
	KeyID    string
	SecretID string
}

func (auth *SigAuth) Sign(r *http.Request) {
	headers := []string{signHeaderRequestTarget, signHeaderDate}
	signer, err := httpsig.NewRequestSigner(auth.KeyID, auth.SecretID, signAlgorithm)
	if err != nil {
		log.Println(err)
		return
	}
	err = signer.SignRequest(r, headers, nil)
	if err != nil {
		log.Println(err)
	}
}

type BasicAuth struct {
	Username string
	Password string
}

func (auth *BasicAuth) Sign(r *http.Request) {
	r.SetBasicAuth(auth.Username, auth.Password)
}

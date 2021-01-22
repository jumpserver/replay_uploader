package service

import (
	"time"

	"jms-upload/pkg/httplib"
)

type option struct {
	// default http://127.0.0.1:8080
	CoreHost string
	TimeOut  time.Duration
	sign     httplib.AuthSign
}

type Option func(*option)

func JMSCoreHost(coreHost string) Option {
	return func(o *option) {
		o.CoreHost = coreHost
	}
}

func JMSTimeOut(t time.Duration) Option {
	return func(o *option) {
		o.TimeOut = t
	}
}
func JMSAuthSign(sign httplib.AuthSign) Option {
	return func(o *option) {
		o.sign = sign
	}
}

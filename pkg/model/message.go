package model

import (
	"encoding/json"
)

const (
	CodeOK  = 200
	CodeErr = 400
)

type Msg struct {
	Err  string `json:"err,omitempty"`
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func (m Msg) String() string {
	result, _ := json.Marshal(m)
	return string(result)
}


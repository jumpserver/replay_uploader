package upload

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReturnMsg(writer io.Writer, m Msg) {
	_, _ = fmt.Fprint(writer, m)
}

func ReturnErrorMsg(msg string, err error) {
	m := Msg{
		Err:  err.Error(),
		Msg:  msg,
		Code: CodeErr,
	}
	ReturnMsg(os.Stderr, m)
}

func ReturnSuccessMsg(msg string) {
	m := Msg{
		Msg:  msg,
		Code: CodeOK,
	}
	ReturnMsg(os.Stdout, m)
}

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

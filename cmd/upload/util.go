package upload

import (
	"fmt"
	"io"
	"os"

	"jms-upload/pkg/model"
)

func ReturnMsg(writer io.Writer, m model.Msg) {
	_, _ = fmt.Fprint(writer, m)
}

func ReturnErrorMsg(msg string, err error) {
	m := model.Msg{
		Err:  err.Error(),
		Msg:  msg,
		Code: model.CodeErr,
	}
	ReturnMsg(os.Stderr, m)
}

func ReturnSuccessMsg(msg string) {
	m := model.Msg{
		Msg:  msg,
		Code: model.CodeOK,
	}
	ReturnMsg(os.Stdout, m)
}

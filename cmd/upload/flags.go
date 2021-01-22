package upload

import (
	"flag"
)

var (
	coreHost   string
	accessKey  string
	replayPath string
	targetDate string
)

func init() {
	flag.StringVar(&coreHost, "url", "http://127.0.0.1:8080", "JumpServer URL，JMS Core的地址")
	flag.StringVar(&accessKey, "key", "", "Access key，组件使用的认证Key")
	flag.StringVar(&replayPath, "file", "", "Replay file path，录像文件路径")
	flag.StringVar(&targetDate, "date", "", "Target date, 默认当前日期，格式 2021-01-20")

}

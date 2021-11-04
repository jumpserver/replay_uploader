package upload

import (
	"flag"
)

/*
	CoreHost 		指 JumpServer 的地址 通常是 http://127.0.0.1:8080
	accessKey 		指 经过 Base64 的 access_key 格式 keyId:secretId
	accessKeyFile   指 保存了access key 的文件，内容格式 keyId:secretId
	replayPath		指 录像文件，未经gzip压缩的录像文件
	targetDate		指 会话生成日期，上传 存储时需要
	sid				指 会话的ID
	forceDelete		指 成功上传之后，是否删除录像 file 源文件
*/

var (
	coreHost      string
	accessKey     string
	accessKeyFile string
	replayPath    string
	targetDate    string
	sid           string
	forceDelete   bool

	infoFlag   = false
)

func init() {
	flag.StringVar(&coreHost, "url", "http://127.0.0.1:8080", "JumpServer URL，JMS Core的地址")
	flag.StringVar(&accessKey, "key", "", "Access key，组件使用的认证Key")
	flag.StringVar(&accessKeyFile, "keyfile", "", "Key file，存储 access key 的文件路径")
	flag.StringVar(&replayPath, "file", "", "Replay file path，录像文件路径")
	flag.StringVar(&sid, "sid", "", "Session ID, 会话ID")
	flag.StringVar(&targetDate, "date", "", "Target date, 默认当前日期，格式 2021-01-20")
	flag.BoolVar(&forceDelete, "remove", false, "成功上传后，是否删除原 file 文件")

	flag.BoolVar(&infoFlag, "V", false, "version info")
}

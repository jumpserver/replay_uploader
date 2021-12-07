package common

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/jumpserver/replay_uploader/jms-sdk-go/model"
	"github.com/jumpserver/replay_uploader/util"
)

type ReplayFile struct {
	ID          string
	TargetDate  string
	AbsFilePath string
	Version     model.ReplayVersion
}

func (r *ReplayFile) TargetPath() string {
	gzFilename := r.GetGzFilename()
	return strings.Join([]string{r.TargetDate, gzFilename}, "/")
}

func (r *ReplayFile) GetGzFilename() string {
	suffixGz := ".replay.gz"
	switch r.Version {
	case model.Version3:
		suffixGz = ".cast.gz"
	case model.Version2:
		suffixGz = ".replay.gz"
	}
	return r.ID + suffixGz
}

/*
koko   文件名为 sid | sid.replay.gz | sid.cast | sid.cast.gz
lion   文件名为 sid | sid.replay.gz
omnidb 文件名为 sid.cast | sid.cast.gz
xrdp   文件名为 sid.guac

如果存在日期目录，targetDate 使用日期目录的

	文件路径名称中解析 录像文件信息
*/

const suffixGuac = ".guac"

var suffixesMap = map[string]model.ReplayVersion{
	suffixGuac:           model.Version2,
	model.SuffixCast:     model.Version3,
	model.SuffixCastGz:   model.Version3,
	model.SuffixReplayGz: model.Version2}

func ParseSessionID(replayFilePath string) (string, bool) {
	filename := filepath.Base(replayFilePath)
	if len(filename) == 36 && util.IsUUID(filename) {
		return filename, true
	}
	sid := strings.Split(filename, ".")[0]
	if !util.IsUUID(sid) {
		return "", false
	}
	return sid, true
}

func ParseReplayVersion(filename string) (model.ReplayVersion, bool) {
	for suffix := range suffixesMap {
		if strings.HasSuffix(filename, suffix) {
			return suffixesMap[suffix], true

		}
	}
	return model.UnKnown, false
}

func ParseDateFromPath(replayFilePath string) (string, bool) {
	dirPath := filepath.Dir(replayFilePath)
	dirName := filepath.Base(dirPath)
	if _, err := time.Parse("2006-01-02", dirName); err == nil {
		return dirName, true
	}
	return "", false
}

package upload

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jumpserver/replay_uploader/jms-sdk-go/model"
	"github.com/jumpserver/replay_uploader/util"
)

/*
koko   文件名为 sid | sid.replay.gz | sid.cast | sid.cast.gz
lion   文件名为 sid | sid.replay.gz
omnidb 文件名为 sid.cast | sid.cast.gz
xrdp   文件名为 sid.guac

如果存在日期目录，targetDate 使用日期目录的
*/

type ReplayFile struct {
	ID        string
	AbsPath   string
	AbsGzPath string
	IsGzip    bool

	Version model.ReplayVersion
}

func ScanFromDirPath(dirPath string) ([]ReplayFile, error) {
	allFiles := make([]ReplayFile, 0, 100)
	walkFunc := func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if replayFile, ok := parseReplayFilename(info.Name()); ok {
			allFiles = append(allFiles, replayFile)
		}
		return nil
	}
	if err := filepath.Walk(dirPath, walkFunc); err != nil {
		return nil, err
	}
	return allFiles, nil
}

func parseReplayFilename(filename string) (replay ReplayFile, ok bool) {
	// 未压缩的旧录像文件名格式是一个 UUID
	if len(filename) == 36 && util.IsUUID(filename) {
		replay.ID = filename
		replay.Version = model.Version2
		ok = true
		return
	}
	if replay.ID, replay.Version, ok = isReplayFile(filename); ok {
		replay.IsGzip = isGzipFile(filename)
	}
	return
}

const suffixGuac = ".guac"

var suffixesMap = map[string]model.ReplayVersion{
	suffixGuac:           model.Version2,
	model.SuffixCast:     model.Version3,
	model.SuffixCastGz:   model.Version3,
	model.SuffixReplayGz: model.Version2}

func isReplayFile(filename string) (id string, version model.ReplayVersion, ok bool) {
	for suffix := range suffixesMap {
		if strings.HasSuffix(filename, suffix) {
			sidName := strings.Split(filename, ".")[0]
			if util.IsUUID(sidName) {
				id = sidName
				version = suffixesMap[suffix]
				ok = true
				return
			}
		}
	}
	return
}

func isGzipFile(filename string) bool {
	return strings.HasSuffix(filename, model.SuffixGz)
}

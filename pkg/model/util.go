package model

import (
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 录像的文件必须是固定格式 sessionID.replay.gz; eg: 90b11402-39df-4ba9-b14a-2344b8585888.replay.gz
func ParseSessionID(replayFilePath string) (string, error) {
	fileName := filepath.Base(replayFilePath)
	sid := strings.TrimSuffix(fileName, SuffixReplayFileName)
	if _, err := uuid.FromString(sid); err != nil {
		return "", err
	}
	return sid, nil
}

func GetCurrentDate() string {
	return time.Now().Format("2006-01-02")
}

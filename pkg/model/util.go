package model

import (
	"compress/gzip"
	"encoding/base64"
	"io"
	"os"
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

func IsValidateSessionID(sid string) bool {
	_, err := uuid.FromString(sid)
	return err == nil
}

func DecodeBase64String(p string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(p)
	return string(result), err
}

func CompressToGzipFile(srcPath, dstPath string) error {
	sf, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer sf.Close()
	df, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer df.Close()
	writer := gzip.NewWriter(df)
	writer.Name = dstPath
	writer.ModTime = time.Now().UTC()
	_, err = io.Copy(writer, sf)
	if err != nil {
		return err
	}
	return writer.Close()
}

func IsGzipFile(gzipFile string) bool {
	extensionName := filepath.Ext(gzipFile)
	return strings.HasSuffix(extensionName, ".gz")
}

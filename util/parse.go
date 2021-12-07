package util

import (
	"fmt"
	"path/filepath"
	"strings"
)

func ParseSessionID(replayFilePath string) (string, error) {
	fileName := filepath.Base(replayFilePath)
	if IsUUID(fileName) {
		return fileName, nil
	}
	filenameSlice := strings.Split(fileName, ".")
	if IsUUID(filenameSlice[0]) {
		return filenameSlice[0], nil
	}
	return "", fmt.Errorf("do not contains session sid %s", replayFilePath)
}

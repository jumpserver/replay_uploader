package model

import (
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func ParseSessionID(replayFilePath string) (string, error) {
	fileName := filepath.Base(replayFilePath)
	if IsValidateSessionID(fileName) {
		return fileName, nil
	}
	filenameSlice := strings.Split(fileName, ".")
	if IsValidateSessionID(filenameSlice[0]) {
		return filenameSlice[0], nil
	}
	return "", fmt.Errorf("do not contains session sid %s", replayFilePath)
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
	writer.Name = filepath.Base(srcPath)
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

func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}

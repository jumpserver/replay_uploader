package util

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"time"
)

func CurrentDate() string {
	return time.Now().Format("2006-01-02")
}

func DecodeBase64String(p string) (string, error) {
	result, err := base64.StdEncoding.DecodeString(p)
	return string(result), err
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

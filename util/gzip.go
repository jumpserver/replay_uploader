package util

import (
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

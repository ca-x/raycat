package readkit

import (
	"bytes"
	"os"
	"path/filepath"
	"raycat/internal/pkg/tinypool"
)

var contentPool = tinypool.New[bytes.Buffer](tinypool.BufReset)

// ReadAll read all data to []byte in a path
func ReadAll(path string) ([]byte, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if !fileInfo.IsDir() {
		return os.ReadFile(path)
	}

	return readDirRecursively(path)
}

// readDirRecursively reads a directory and its subdirectories recursively
func readDirRecursively(dirPath string) ([]byte, error) {
	allContent := contentPool.Get()
	defer contentPool.Free(allContent)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())
		if entry.IsDir() {
			subContent, err := readDirRecursively(fullPath)
			if err != nil {
				return nil, err
			}
			allContent.Write(subContent)
		} else {
			if entry.Type().IsRegular() {
				f, err := os.Open(fullPath)
				if err != nil {
					continue
				}
				_, err = f.WriteTo(allContent)
				if err != nil {
					f.Close()
					continue
				}
				f.Close()
			}
		}
	}
	return allContent.Bytes(), nil
}

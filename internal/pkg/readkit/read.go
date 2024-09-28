package readkit

import (
	"os"
	"path/filepath"
)

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
	var allContent []byte
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
			allContent = append(allContent, subContent...)
		} else {
			if entry.Type().IsRegular() {
				content, err := os.ReadFile(fullPath)
				if err != nil {
					return nil, err
				}
				allContent = append(allContent, content...)
			}
		}
	}
	return allContent, nil
}

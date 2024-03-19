package fileutil

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// ScanFiles scans files in a folder with specified file extensions and returns their paths.
func ScanFiles(directory string, extensions []string) ([]string, error) {
	var files []string

	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// Check if the file extension is in the whitelist
			ext := filepath.Ext(path)
			for _, allowedExt := range extensions {
				if strings.EqualFold(ext, allowedExt) {
					files = append(files, path)
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

// ReadFileLines reads the lines of a file and returns them as a slice of strings.
func ReadFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

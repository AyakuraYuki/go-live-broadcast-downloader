package file

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// IsPathExist returns if given path (file or dir) exists
func IsPathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// SaveFile saves data to file
func SaveFile(path string, value string) {
	f, err := os.Create(path)
	defer f.Close()
	if err != nil {
		return
	}

	buf := bufio.NewWriter(f)
	_, _ = fmt.Fprintln(buf, value)
	_ = buf.Flush()
}

// ReadFile reads the whole file and result to a string
func ReadFile(filePth string) string {
	b := ReadFileByte(filePth)
	if b == nil {
		return ""
	}
	return string(b)
}

// ReadFileByte reads the whole file and result to a bytes
func ReadFileByte(filePth string) []byte {
	f, err := os.Open(filePth)
	defer f.Close()
	if err != nil {
		return nil
	}

	r, err := io.ReadAll(f)
	return r
}

// ReadLines returns content by lines in file
func ReadLines(filename string) []string {
	lines := make([]string, 0)
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		return lines
	}
	reader := bufio.NewReader(f)

	for {
		line, err0 := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
		if err0 == io.EOF {
			break
		}
	}

	return lines
}

// WriteLines writes file by lines
func WriteLines(filename string, value []string) {
	f, err := os.Create(filename)
	defer f.Close()
	if err != nil {
		return
	}
	writer := bufio.NewWriter(f)
	for _, v := range value {
		_, _ = fmt.Fprintln(writer, v)
	}
	_ = writer.Flush()
}

// ListDir scans the given dirPath (excludes sub dirs), results to a fileList
func ListDir(dirPath string) (fileList []string, err error) {
	fileList = make([]string, 0)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fileList, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileList = append(fileList, filepath.Join(dirPath, entry.Name()))
	}
	return fileList, nil
}

// WalkDir scans the given dirPath (includes sub dirs), save filenames to *fileList
func WalkDir(dirPath string, fileList *[]string) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			WalkDir(filepath.Join(dirPath, entry.Name()), fileList)
		} else {
			*fileList = append(*fileList, filepath.Join(dirPath, entry.Name()))
		}
	}
}

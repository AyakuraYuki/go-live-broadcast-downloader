package file

import (
	"bufio"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// ReadBytes reads the whole file and returns the content bytearray.
func ReadBytes(path string) []byte {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	bs, _ := io.ReadAll(f)
	return bs
}

// ReadString reads the whole file and returns the text content.
func ReadString(path string) string {
	bs := ReadBytes(path)
	if len(bs) == 0 {
		return ""
	}
	return string(bs)
}

// ReadLines returns the text content by lines
func ReadLines(path string) (lines []string) {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer func(f *os.File) { _ = f.Close() }(f)
	lines = make([]string, 0)
	reader := bufio.NewReader(f)
	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line != "" {
			lines = append(lines, line)
		}
		if err == io.EOF {
			break
		}
	}
	return lines
}

// ListDir scans the given path without recurse, results to a list of file abs paths
func ListDir(dirPath string) (files []string, err error) {
	files = make([]string, 0)
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return files, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		files = append(files, filepath.Join(dirPath, entry.Name()))
	}
	return files, nil
}

// Walk scans the given path with recurse, results to a file list in absolute-path
func Walk(base string) (files []string, err error) {
	files = make([]string, 0)
	err = filepath.Walk(base, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

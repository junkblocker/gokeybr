package fs

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const FileAccess = 0o644

func homeFilePath(name string) string {
	dir := os.Getenv("XDG_STATE_HOME")
	if dir == "" {
		dir = filepath.Join(
			os.Getenv("HOME"),
			".local",
			"state")
	}
	return filepath.Join(dir, "gokeybr", name)
}

func mkdir() {
	dir := homeFilePath("/")
	if _, err := os.Stat(dir); err != nil {
		_ = os.MkdirAll(dir, os.ModePerm)
	}
}

func SaveJSON(filename string, o interface{}) error {
	data, err := json.MarshalIndent(o, "", " ")
	if err != nil {
		return err
	}
	mkdir()
	return os.WriteFile(homeFilePath(filename), data, FileAccess)
}

func LoadJSON(filename string, v interface{}) error {
	mkdir()
	data, err := os.ReadFile(homeFilePath(filename))
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func AppendJSONLine(filename string, v interface{}) (err error) {
	mkdir()
	var f *os.File
	if f, err = os.OpenFile(homeFilePath(filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, FileAccess); err != nil {
		return err
	}
	defer f.Close()
	var data []byte
	if data, err = json.Marshal(v); err != nil {
		return err
	}
	_, err = fmt.Fprintln(f, string(data))
	return err
}

type JSONLinesIterator struct {
	scanner *bufio.Scanner
	file    *os.File
}

func NewJSONLinesIterator(filename string) (*JSONLinesIterator, error) {
	file, err := os.Open(homeFilePath(filename))
	if err != nil {
		return nil, err
	}
	return &JSONLinesIterator{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (i JSONLinesIterator) Close() {
	_ = i.file.Close()
}

func (i JSONLinesIterator) UnmarshalNextLine(v interface{}) (bool, error) {
	if !i.scanner.Scan() {
		return false, i.scanner.Err()
	}

	return true, json.Unmarshal(i.scanner.Bytes(), v)
}

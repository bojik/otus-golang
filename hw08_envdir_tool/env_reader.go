package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

type Environment map[string]EnvValue

var (
	ErrDirDoesNotExist = errors.New("directory does not exist")
	ErrIsNotDir        = errors.New("given path is not dir")
)

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	stat, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("%s: %w", dir, ErrDirDoesNotExist)
	}
	if !stat.IsDir() {
		return nil, fmt.Errorf("%s: %w", dir, ErrIsNotDir)
	}
	files, err1 := os.ReadDir(dir)
	if err1 != nil {
		return nil, fmt.Errorf("%s: %w", dir, err1)
	}
	envs := Environment{}
	for _, file := range files {
		if strings.Contains(file.Name(), "=") {
			continue
		}
		fullName := fmt.Sprintf("%s%c%s",
			strings.TrimRight(dir, string(os.PathSeparator)),
			os.PathSeparator,
			file.Name(),
		)
		stat, err := os.Stat(fullName)
		if err != nil {
			return nil, err
		}
		if stat.IsDir() {
			continue
		}
		if stat.Size() == 0 {
			envs[file.Name()] = EnvValue{
				NeedRemove: true,
			}
			continue
		}
		fl, err := readFirstLine(fullName)
		if err != nil {
			return nil, err
		}
		envs[file.Name()] = EnvValue{
			Value: clearString(fl),
		}
	}
	return envs, nil
}

// readFirstLine reads first line only from file given by fileName.
func readFirstLine(fileName string) (string, error) {
	fp, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer fp.Close()
	con, err := io.ReadAll(fp)
	if err != nil {
		return "", err
	}
	return strings.Split(string(con), "\n")[0], nil
}

// clearString prepares string from usage.
func clearString(s string) string {
	ret := strings.ReplaceAll(s, "\x00", "\n")
	re := regexp.MustCompile("[\t ].+$")
	ret = string(re.ReplaceAll([]byte(ret), []byte{}))
	return ret
}

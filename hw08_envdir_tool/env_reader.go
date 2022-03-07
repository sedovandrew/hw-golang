package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// readFirstLineFromFile reads the first line from a file.
func readFirstLineFromFile(filename string) (string, error) {
	file, err := os.Open(filename)
	defer func() {
		err = file.Close()
	}()
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	return firstLine, err
}

// lineEditing changes a line based on certain conditions.
func lineEditing(line string) string {
	line = strings.TrimRight(line, " \t")
	lineBytes := bytes.ReplaceAll([]byte(line), []byte{0x00}, []byte("\n"))
	return string(lineBytes)
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileInfoSlice, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, fileInfo := range fileInfoSlice {
		// Skip directories.
		if fileInfo.IsDir() {
			continue
		}

		filename := fileInfo.Name()
		// Skip files with "=" in filename.
		if strings.Contains(filename, "=") {
			continue
		}
		firstLine, err := readFirstLineFromFile(path.Join(dir, filename))
		if err != nil {
			continue
		}
		value := EnvValue{
			Value: lineEditing(firstLine),
		}
		if firstLine == "" {
			value.NeedRemove = true
		}
		env[filename] = value
	}
	return env, nil
}

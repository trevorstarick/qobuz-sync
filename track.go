package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pkg/errors"
)

type Track struct {
	cache map[string]string
	f     *os.File

	Path string `json:"path"`
}

func NewTracker(path string) (*Track, error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open file")
	}

	t := &Track{
		cache: make(map[string]string),
		f:     f,
		Path:  path,
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read file")
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		if line == "" {
			continue
		}

		delim := ": "

		index := strings.Index(line, delim)
		if index == -1 {
			log.Println(warn, "invalid line in tracker file:", line)

			continue
		}

		id := line[:index]
		path := line[index+len(delim):]

		_, err := os.Stat(path)
		if err != nil {
			log.Println(warn, "unable to stat file:", path)

			continue
		}

		t.cache[id] = path
	}

	return t, nil
}

func (t *Track) Set(key string, value string) error {
	if _, ok := t.cache[key]; ok {
		return nil
	}

	t.cache[key] = value
	_, err := t.f.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	if err != nil {
		return errors.Wrap(err, "unable to write to file")
	}

	return nil
}

func (t *Track) Get(key string) (string, error) {
	if v, ok := t.cache[key]; ok {
		return v, nil
	}

	return "", errors.New("key not found")
}

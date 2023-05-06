package client

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
)

type Tracker struct {
	cache map[string]string
	file  *os.File

	Path string `json:"path"`
}

func NewTracker(path string) (*Tracker, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, common.FilePerm)
	if err != nil {
		return nil, errors.Wrap(err, "unable to open file")
	}

	tracker := &Tracker{
		cache: make(map[string]string),
		file:  file,
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
			log.Warn().Msgf("invalid line in tracker file: %v", line)

			continue
		}

		trackOrAlbumID := line[:index]
		path := line[index+len(delim):]

		_, err := os.Stat(path)
		if err != nil {
			log.Warn().Msgf("unable to stat file: %v", path)

			continue
		}

		tracker.cache[trackOrAlbumID] = path
	}

	return tracker, nil
}

func (tracker *Tracker) Set(key string, value string) error {
	if _, ok := tracker.cache[key]; ok {
		return nil
	}

	tracker.cache[key] = value

	_, err := tracker.file.WriteString(fmt.Sprintf("%s: %s\n", key, value))
	if err != nil {
		return errors.Wrap(err, "unable to write to file")
	}

	return nil
}

func (tracker *Tracker) Get(key string) (string, error) {
	if v, ok := tracker.cache[key]; ok {
		return v, nil
	}

	return "", errors.New("key not found")
}

func (tracker *Tracker) Close() error {
	err := tracker.file.Close()
	if err != nil {
		return errors.Wrap(err, "unable to close file")
	}

	return nil
}

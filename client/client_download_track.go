package client

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
)

//nolint:cyclop // TODO: refactor
func (client *Client) downloadFile(trackID, path string) error {
	url, err := client.TrackGetFileURL(trackID, QualityMAX)
	if err != nil {
		return errors.Wrap(err, "failed to get track file url")
	}

	if url.MimeType != "audio/flac" {
		log.Warn().Msgf("not FLAC got %v: %v", url.MimeType, path)
		path = strings.ReplaceAll(path, ".flac", ".mp3")
	}

	if _, err := os.Stat(strings.TrimSuffix(path, ".part")); err == nil {
		return errors.Wrap(common.ErrAlreadyExists, "cached")
	}

	// do some basic verification that the url is valid
	if !strings.HasPrefix(url.URL, "https://streaming-qobuz-std.akamaized.net/file?") {
		return errors.New("was given an invalid streaming url from qobuz")
	}

	res, err := http.Get(url.URL) //nolint:noctx
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = errors.Wrap(err, "failed to close m3u file")
		}
	}()

	audioFile, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, common.FilePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	defer func() {
		if syncErr := audioFile.Sync(); syncErr != nil {
			err = errors.Wrap(err, "failed to sync file")
		}

		if closeErr := audioFile.Close(); closeErr != nil {
			err = errors.Wrap(err, "failed to close file")
		}
	}()

	_, err = io.Copy(audioFile, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
}

func (client *Client) downloadFileAndSetMetadata(trackID, path string, metadata common.Metadata) error {
	partialPath := path + ".part"

	err := client.downloadFile(trackID, partialPath)
	if err != nil {
		return errors.Wrap(err, "failed to download track")
	}

	err = SetTags(partialPath, metadata)
	if err != nil {
		return errors.Wrap(err, "failed to set tags")
	}

	return nil
}

func (client *Client) downloadTrack(trackID string) error {
	_, err := client.trackTracker.Get(trackID)
	if err == nil {
		return errors.Wrap(common.ErrAlreadyExists, "cached")
	}

	track, err := client.TrackGet(trackID)
	if err != nil {
		return errors.Wrap(err, "failed to get track")
	}

	trackPath := filepath.Join(client.baseDir, track.Path())

	err = os.MkdirAll(filepath.Dir(trackPath), common.DirPerm)
	if err != nil {
		return errors.Wrap(err, "failed to create directory")
	}

	_, err = os.Stat(trackPath)
	if err == nil {
		err = client.trackTracker.Set(trackID, trackPath)
		if err != nil {
			return errors.Wrap(err, "failed to set track as downloaded")
		}

		return common.ErrAlreadyExists
	}

	err = client.downloadFileAndSetMetadata(trackID, trackPath, track.Metadata())
	if err != nil {
		return errors.Wrap(err, "failed to download and set metadata")
	}

	log.Info().Msgf("downloaded track: %v", trackPath)

	err = client.trackTracker.Set(trackID, trackPath)
	if err != nil {
		return errors.Wrap(err, "failed to set track as downloaded")
	}

	return nil
}

func (client *Client) DownloadTrack(trackID string) error {
	err := client.downloadTrack(trackID)
	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			dir, _ := client.trackTracker.Get(trackID)
			log.Info().Msgf("track already downloaded: %v", dir)

			return nil
		}

		return errors.Wrap(err, "failed to download track")
	}

	return nil
}

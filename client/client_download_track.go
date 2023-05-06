package client

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
)

func (client *Client) download(trackID string, path string) error {
	url, err := client.TrackGetFileURL(trackID, QualityMAX)
	if err != nil {
		return errors.Wrap(err, "failed to get track file url")
	}

	if url.FormatID == 0 {
		return common.ErrUnavailable
	}

	if url.MimeType != "audio/flac" {
		log.Warn().Msgf("not FLAC got %v: %v", url.MimeType, path)
		path = strings.ReplaceAll(path, ".flac", ".mp3")
	}

	if _, err := os.Stat(strings.TrimSuffix(path, ".part")); err == nil {
		return errors.Wrap(common.ErrAlreadyExists, "cached")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.URL, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, common.FilePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
}

func (client *Client) downloadAndSetMetadata(trackID string, path string, metadata common.Metadata) error {
	partialPath := path + ".part"

	err := client.download(trackID, partialPath)
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

	err = client.downloadAndSetMetadata(trackID, trackPath, track.Metadata())
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
	return client.downloadTrack(trackID)
}

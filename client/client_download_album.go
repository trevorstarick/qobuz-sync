package client

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/responses"
)

func (client *Client) downloadAlbumTrack(track *responses.Track, album *responses.Album) error {
	if track.Album == nil {
		track.Album = album
	}

	_, err := os.Stat(track.Path())
	if err == nil {
		log.Info().Msgf("track already exists, skipping: %v", track.Path())

		return nil
	}

	err = client.downloadTrack(strconv.Itoa(track.ID))
	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			log.Info().Msgf("track already exists, skipping: %v", track.Path())

			return nil
		}

		return err
	}

	return nil
}

func (client *Client) downloadAlbum(albumID string) (*responses.Album, error) {
	if !client.force {
		_, err := client.albumTracker.Get(albumID)
		if err == nil {
			return nil, errors.Wrap(common.ErrAlreadyExists, "cached")
		}
	}

	album, err := client.AlbumGet(albumID)
	if err != nil {
		return nil, errors.Wrap(err, "album get")
	}

	albumDir := filepath.Join(client.baseDir, album.Path())

	err = os.MkdirAll(albumDir, common.DirPerm)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create directory")
	}

	for i := range album.Tracks.Items {
		err = client.downloadAlbumTrack(&album.Tracks.Items[i], album.Album)
		if err != nil {
			path := album.Tracks.Items[i].Path()
			log.Error().Err(err).Msgf("failed to download track, skipping: %v", path)

			continue
		}
	}

	if album.ReleasedAt < int(time.Now().Unix()) {
		err = client.albumTracker.Set(albumID, albumDir)
		if err != nil {
			return nil, errors.Wrap(err, "failed to set album as downloaded")
		}
	} else {
		log.Info().Msgf("album not released yet, not setting as downloaded: %v", album.Path())
	}

	return album.Album, nil
}

func (client *Client) DownloadAlbum(albumID string) error {
	album, err := client.downloadAlbum(albumID)
	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			dir, _ := client.albumTracker.Get(albumID)
			log.Info().Msgf("album already downloaded: %v", dir)

			return nil
		} else if errors.Is(err, common.ErrNotFound) {
			log.Warn().Msgf("album not found: %v", albumID)

			return nil
		}

		return errors.Wrap(err, "failed to download album")
	}

	albumDir := filepath.Join(client.baseDir, album.Path())

	err = album.DownloadAlbumArt(albumDir)
	if err != nil {
		if errors.Is(err, common.ErrAlreadyExists) {
			log.Info().Msgf("album art already exists, skipping: %v/album.jpg", album.Path())
		} else {
			log.Warn().Msgf("failed to download album art, skipping: %v", err)
		}
	}

	log.Info().Msgf("downloaded album: %v", albumDir)

	return nil
}

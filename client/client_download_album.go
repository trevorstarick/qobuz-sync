package client

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/responses"
)

func (client *Client) downloadAlbum(albumID string) (*responses.Album, error) {
	_, err := client.albumTracker.Get(albumID)
	if err == nil {
		return nil, errors.Wrap(common.ErrAlreadyExists, "cached")
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

	for _, track := range album.Tracks.Items {
		if track.Album == nil {
			track.Album = album.Album
		}

		_, err = os.Stat(track.Path())
		if err == nil {
			log.Info().Msgf("track already exists, skipping: %v", track.Path())

			continue
		}

		err = client.downloadTrack(strconv.Itoa(track.ID))
		if err != nil {
			if errors.Is(err, common.ErrAlreadyExists) {
				log.Info().Msgf("track already exists, skipping: %v", track.Path())
			} else {
				log.Warn().Msgf("failed to download track, skipping: %v", err)
			}

			continue
		}
	}

	err = client.albumTracker.Set(albumID, albumDir)
	if err != nil {
		return nil, errors.Wrap(err, "failed to set album as downloaded")
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

package client

import (
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
)

func (*Client) Search(_ string) error {
	return common.ErrNotImplemented
}

func (*Client) GetArtist(_ string) error {
	return common.ErrNotImplemented
}

func (client *Client) FavoriteAlbums() error {
	offset := 0

	for {
		res, err := client.FavoriteGetUserFavorites(ListTypeALBUM, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for i := range res.Albums.Items {
			album := &res.Albums.Items[i]

			err = client.DownloadAlbum(album.ID)
			if err != nil {
				if errors.Is(err, common.ErrAlreadyExists) {
					dir, _ := client.albumTracker.Get(album.ID)
					log.Info().Msgf("album already exists: %v", dir)
				} else {
					log.Warn().Msgf("unable to download album, skipping: %v", err)
				}

				continue
			}
		}

		if res.Albums.Offset+res.Albums.Limit >= res.Albums.Total {
			break
		}

		offset += res.Albums.Limit
	}

	return nil
}

func (client *Client) FavoriteTracks() error {
	offset := 0

	for {
		res, err := client.FavoriteGetUserFavorites(ListTypeTRACK, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for _, track := range res.Tracks.Items {
			err = client.DownloadTrack(strconv.Itoa(track.ID))
			if err != nil {
				if errors.Is(err, common.ErrAlreadyExists) {
					path, _ := client.trackTracker.Get(strconv.Itoa(track.ID))
					log.Info().Msgf("track already exists: %v", path)
				} else {
					log.Warn().Msgf("unable to download track, skipping: %v", err)
				}

				continue
			}

			albumDir := filepath.Join(client.baseDir, track.Album.Path())

			err = track.Album.DownloadAlbumArt(albumDir)
			if err != nil {
				if errors.Is(err, common.ErrAlreadyExists) {
					log.Info().Msgf("album art already exists: %v/album.jpg", albumDir)
				} else {
					log.Warn().Msgf("unable to download album art, skipping: %v", err)
				}

				continue
			}
		}

		if res.Tracks.Offset+res.Tracks.Limit >= res.Tracks.Total {
			break
		}

		offset += res.Tracks.Limit
	}

	return nil
}

// https://open.qobuz.com/track/6451477
// https://open.qobuz.com/album/0603497932191
// https://open.qobuz.com/artist/34527
// https://open.qobuz.com/playlist/2418316
func (client *Client) Link(link string) error {
	u, err := url.Parse(link) //nolint:varnamelen
	if err != nil {
		return errors.Wrap(err, "unable to parse link")
	}

	switch u.Host {
	case "open.qobuz.com":
		parts := strings.Split(u.Path, "/")

		if len(parts) < 3 { //nolint:gomnd
			return errors.Wrap(common.ErrNotImplemented, "unsupported link")
		}

		parts = parts[1:]

		switch parts[0] {
		case "track":
			return client.DownloadTrack(parts[1])
		case "album":
			return client.DownloadAlbum(parts[1])
		case "artist":
			return common.ErrNotImplemented
		case "playlist":
			return client.DownloadPlaylist(parts[1])
		default:
			return errors.Wrap(common.ErrNotImplemented, "unsupported link")
		}
	default:
		return errors.Wrap(common.ErrNotImplemented, "unsupported host")
	}
}

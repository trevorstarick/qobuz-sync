package client

import (
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/helpers"
)

func (client *Client) DownloadPlaylist(playlistID string) error {
	res, err := client.PlaylistGet(playlistID)
	if err != nil {
		return errors.Wrap(err, "playlist get")
	}

	playlistDir := filepath.Join(client.baseDir, "_playlist", helpers.SanitizeStringToPath(res.Name))

	err = os.MkdirAll(playlistDir, common.DirPerm)
	if err != nil {
		return errors.Wrap(err, "failed to create playlist dir")
	}

	m3uFile, err := os.Create(filepath.Join(playlistDir, "playlist.m3u"))
	if err != nil {
		return errors.Wrap(err, "failed to create m3u file")
	}

	defer m3uFile.Close()

	_, _ = m3uFile.WriteString("#EXTM3U\n")
	_, _ = m3uFile.WriteString("#EXTENC: UTF-8\n")
	_, _ = m3uFile.WriteString("#PLAYLIST: " + res.Name + "\n")
	_, _ = m3uFile.WriteString("#EXTID: " + playlistID + "\n")

	for _, track := range res.Tracks.Items {
		err = client.DownloadTrack(strconv.Itoa(track.ID))
		if err != nil {
			if errors.Is(err, common.ErrAlreadyExists) {
				path, _ := client.trackTracker.Get(strconv.Itoa(track.ID))
				log.Info().Msgf("track already exists, skipping: %v", path)

				_, _ = m3uFile.WriteString("#EXTINF:" + strconv.Itoa(track.Duration) + "," + track.Album.Artist.Name + " - " + track.Title + "\n") //nolint:lll // m3u format

				relPath, err := filepath.Rel(playlistDir, path)
				if err != nil {
					return errors.Wrap(err, "failed to get relative path")
				}

				_, _ = m3uFile.WriteString(relPath + "\n")
			} else {
				log.Warn().Err(err).Msgf("failed to download track, skipping: %v", strconv.Itoa(track.ID))
			}
		}
	}

	log.Info().Msgf("downloaded playlist: %v", playlistDir)

	return nil
}

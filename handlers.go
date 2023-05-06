package main

import (
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	client  *Client
	baseDir string
}

func NewHandler(client *Client, baseDir string) *Handler {
	return &Handler{
		client:  client,
		baseDir: baseDir,
	}
}

func (h *Handler) Album(id string) error {
	dir, err := downloadAlbum(h.client, id, h.baseDir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Info().Msgf("album already exists: %v", dir)
		} else {
			return errors.Wrap(err, "unable to download album")
		}
	}

	return nil
}

func (h *Handler) Track(id string) error {
	res, err := h.client.TrackGet(id)
	if err != nil {
		return errors.Wrap(err, "unable to download track")
	}

	artist := sanitizeStringToPath(res.Album.Artist.Name)
	albumName := sanitizeStringToPath(res.Album.Title)
	dir := filepath.Join(h.baseDir, artist, albumName)

	path, err := downloadTrack(h.client, strconv.Itoa(res.ID), dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Info().Msgf("track already exists: %v", path)
		} else {
			return errors.Wrap(err, "unable to download track")
		}
	}

	err = downloadAlbumArt(res.Album.Image.Large, dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Info().Msgf("album art already exists: %v", dir)
		} else {
			return errors.Wrap(err, "unable to download album art")
		}
	}

	return nil
}

func (h *Handler) FavoriteAlbums() error {
	offset := 0

	for {
		res, err := h.client.FavoriteGetUserFavorites(ListTypeALBUM, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for _, album := range res.Albums.Items {
			dir, err := downloadAlbum(h.client, album.ID, h.baseDir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
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

func (h *Handler) FavoriteTracks() error {
	offset := 0

	for {
		res, err := h.client.FavoriteGetUserFavorites(ListTypeTRACK, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for _, track := range res.Tracks.Items {
			artist := sanitizeStringToPath(track.Album.Artist.Name)
			albumName := sanitizeStringToPath(track.Album.Title)
			dir := filepath.Join(h.baseDir, artist, albumName)

			path, err := downloadTrack(h.client, strconv.Itoa(track.ID), dir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
					log.Info().Msgf("track already exists: %v", path)
				} else {
					log.Warn().Msgf("unable to download track, skipping: %v", err)
				}

				continue
			}

			err = downloadAlbumArt(track.Album.Image.Large, dir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
					log.Info().Msgf("album art already exists: %v", dir)
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

func (h *Handler) Search(query string) error {
	_ = query

	return ErrNotImplemented
}

func (h *Handler) Playlist(id string) error {
	_ = id

	return ErrNotImplemented
}

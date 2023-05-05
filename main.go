package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/pkg/errors"
)

func handleAlbum(c *Client, id string, baseDir string) error {
	dir, err := downloadAlbum(c, id, baseDir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Println(infoMsg, "album already exists:", dir)
		} else {
			return errors.Wrap(err, "unable to download album")
		}
	}

	return nil
}

func handleTrack(c *Client, id string, baseDir string) error {
	res, err := c.TrackGet(id)
	if err != nil {
		return errors.Wrap(err, "unable to download track")
	}

	artist := sanitizeStringToPath(res.Album.Artist.Name)
	albumName := sanitizeStringToPath(res.Album.Title)
	dir := filepath.Join(baseDir, artist, albumName)

	path, err := downloadTrack(c, strconv.Itoa(res.ID), dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Println(infoMsg, "track already exists:", path)
		} else {
			return errors.Wrap(err, "unable to download track")
		}
	}

	err = downloadAlbumArt(res.Album.Image.Large, dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Println(infoMsg, "album art already exists:", dir)
		} else {
			return errors.Wrap(err, "unable to download album art")
		}
	}

	return nil
}

func handleFavoriteAlbums(c *Client, baseDir string) error {
	offset := 0

	for {
		res, err := c.FavoriteGetUserFavorites(ListTypeALBUM, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for _, album := range res.Albums.Items {
			dir, err := downloadAlbum(c, album.ID, baseDir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
					log.Println(infoMsg, "album already exists:", dir)
				} else {
					log.Println(warnMsg, "unable to download album, skipping:", err)
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

func handleFavoriteTracks(c *Client, baseDir string) error {
	offset := 0

	for {
		res, err := c.FavoriteGetUserFavorites(ListTypeTRACK, offset)
		if err != nil {
			return errors.Wrap(err, "unable to get favorites list")
		}

		for _, track := range res.Tracks.Items {
			artist := sanitizeStringToPath(track.Album.Artist.Name)
			albumName := sanitizeStringToPath(track.Album.Title)
			dir := filepath.Join(baseDir, artist, albumName)

			path, err := downloadTrack(c, strconv.Itoa(track.ID), dir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
					log.Println(infoMsg, "track already exists:", path)
				} else {
					log.Println(warnMsg, "unable to download track, skipping:", err)
				}

				continue
			}

			err = downloadAlbumArt(track.Album.Image.Large, dir)
			if err != nil {
				if errors.Is(err, ErrAlreadyExists) {
					log.Println(infoMsg, "album art already exists:", dir)
				} else {
					log.Println(warnMsg, "unable to download album art, skipping:", err)
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

func handleFavorites(c *Client, mode string, baseDir string) error {
	switch mode {
	case "albums":
		err := handleFavoriteAlbums(c, baseDir)
		if err != nil {
			return err
		}
	case "tracks":
		err := handleFavoriteTracks(c, baseDir)
		if err != nil {
			return err
		}
	case "albums+tracks", "tracks+albums":
		err := handleFavoriteTracks(c, baseDir)
		if err != nil {
			return err
		}

		err = handleFavoriteAlbums(c, baseDir)
		if err != nil {
			return err
		}
	default:
		log.Println("usage: qobuz-sync favorites <albums|tracks|albums+tracks>")

		return ErrInvalidArgs
	}

	return nil
}

func run() error {
	var (
		baseDir  = "downloads"
		username = os.Getenv("QOBUZ_USERNAME")
		password = os.Getenv("QOBUZ_PASSWORD")
	)

	if os.Getenv("QOBUZ_BASEDIR") != "" {
		baseDir = os.Getenv("QOBUZ_BASEDIR")
	}

	if username == "" || password == "" {
		log.Println(errMsg, "QOBUZ_USERNAME and QOBUZ_PASSWORD envvars must be set")

		return errors.Wrap(ErrAuthFailed, "missing credentials")
	}

	err := os.MkdirAll(baseDir, dirPerm)
	if err != nil {
		return errors.Wrap(err, "unable to create base dir")
	}

	trackTracker, err = NewTracker(filepath.Join(baseDir, "tracks.txt"))
	if err != nil {
		return errors.Wrap(err, "unable to create track tracker")
	}
	defer trackTracker.Close()

	albumTracker, err = NewTracker(filepath.Join(baseDir, "albums.txt"))
	if err != nil {
		return errors.Wrap(err, "unable to create album tracker")
	}
	defer albumTracker.Close()

	args := os.Args[1:]

	if len(args) == 0 {
		log.Println("usage: qobuz-sync <album|track|favorites>")

		return errors.Wrap(ErrInvalidArgs, "no command specified")
	}

	c, err := NewClient(username, password)
	if err != nil {
		log.Println(errMsg, "unable to login")

		return ErrAuthFailed
	}

	switch args[0] {
	case "album":
		if len(args) != 2 { //nolint:gomnd
			log.Println("usage: qobuz-sync album <id>")

			return errors.Wrap(ErrInvalidArgs, "invalid number of arguments")
		}

		err = handleAlbum(c, args[1], baseDir)
		if err != nil {
			return err
		}
	case "track":
		if len(args) != 2 { //nolint:gomnd
			log.Println("usage: qobuz-sync track <id>")

			return errors.Wrap(ErrInvalidArgs, "invalid number of arguments")
		}

		err = handleTrack(c, args[1], baseDir)
		if err != nil {
			return err
		}
	case "favorites":
		if len(args) != 2 { //nolint:gomnd
			log.Println("usage: qobuz-sync favorites <album|track>")

			return errors.Wrap(ErrInvalidArgs, "invalid number of arguments")
		}

		err = handleFavorites(c, args[1], baseDir)
		if err != nil {
			return err
		}
	default:
		log.Println("usage: qobuz-sync <album|track|favorites>")

		return errors.Wrap(ErrInvalidArgs, "invalid number of arguments")
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		log.Println(errMsg, err)

		if errors.Is(err, ErrInvalidArgs) {
			os.Exit(1)
		} else if errors.Is(err, ErrAuthFailed) {
			os.Exit(2)
		} else {
			os.Exit(3)
		}
	}
}

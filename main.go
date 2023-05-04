package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

var baseDir = "downloads"

func main() {
	username := os.Getenv("QOBUZ_USERNAME")
	password := os.Getenv("QOBUZ_PASSWORD")
	if os.Getenv("QOBUZ_BASEDIR") != "" {
		baseDir = os.Getenv("QOBUZ_BASEDIR")
	}

	if username == "" || password == "" {
		log.Println(err, "QOBUZ_USERNAME and QOBUZ_PASSWORD envvars must be set")
		os.Exit(1)
	}

	os.Args = os.Args[1:]

	if len(os.Args) == 0 {
		log.Println("usage: qobuz-sync <album|track|favorites>")
		os.Exit(1)
	}

	c, err := NewClient(username, password)
	if err != nil {
		log.Println(err, "unable to login")
		os.Exit(1)
	}

	switch os.Args[0] {
	case "album":
		if len(os.Args) != 2 {
			log.Println("usage: qobuz-sync album <id>")

			return
		}

		err := downloadAlbum(c, os.Args[1], baseDir)
		if err != nil {
			log.Println(err, "unable to download album:", err)
			os.Exit(1)
		}

	case "track":
		if len(os.Args) != 2 {
			log.Println("usage: qobuz-sync track <id>")

			return
		}

		res, err := c.TrackGet(os.Args[1])
		if err != nil {
			log.Println(err, "unable to download track:", err)
			os.Exit(1)
		}

		artist := sanitizeStringToPath(res.Album.Artist.Name)
		albumName := sanitizeStringToPath(res.Album.Title)
		dir := filepath.Join(baseDir, artist, albumName)

		path, err := downloadTrack(c, strconv.Itoa(res.ID), dir)
		if err != nil {
			if errors.Is(err, ErrAlreadyExists) {
				log.Println(info, "track already exists:", path)
			} else {
				log.Println(err, "unable to download track:", err)
				os.Exit(1)
			}
		}

		err = downloadAlbumArt(res.Album.Image.Large, dir)
		if err != nil {
			log.Println(err, "unable to download album art:", err)
			os.Exit(1)
		}

	case "favorites":
		if len(os.Args) != 2 {
			log.Println("usage: qobuz-sync favorites <albums|tracks|albums+tracks>")

			return
		}

		var getTracks, getAlbums bool
		switch os.Args[1] {
		case "albums":
			getAlbums = true
		case "tracks":
			getTracks = true
		case "albums+tracks", "tracks+albums":
			getAlbums = true
			getTracks = true
		default:
			log.Println("usage: qobuz-sync favorites <albums|tracks|albums+tracks>")
		}

		if getTracks {
			res, err := c.FavoriteGetUserFavorites(ListTypeTRACK, 0)
			if err != nil {
				log.Println(err, "unable to get favorites list:", err)
				os.Exit(1)
			}

			for _, track := range res.Tracks.Items {
				artist := sanitizeStringToPath(track.Album.Artist.Name)
				albumName := sanitizeStringToPath(track.Album.Title)
				dir := filepath.Join(baseDir, artist, albumName)

				path, err := downloadTrack(c, strconv.Itoa(track.ID), dir)
				if err != nil {
					if errors.Is(err, ErrAlreadyExists) {
						log.Println(info, "track already exists:", path)
					} else {
						log.Println(warn, "unable to download track, skipping:", err)
					}

					continue
				}

				err = downloadAlbumArt(track.Album.Image.Large, dir)
				if err != nil {
					log.Println(warn, "unable to download album art, skipping:", err)

					continue
				}
			}
		}

		if getAlbums {
			res, err := c.FavoriteGetUserFavorites(ListTypeALBUM, 0)
			if err != nil {
				log.Println(err, "unable to get favorites list:", err)
				os.Exit(1)
			}

			for _, album := range res.Albums.Items {
				err = downloadAlbum(c, album.ID, baseDir)
				if err != nil {
					log.Println(warn, "unable to download album, skipping:", err)

					continue
				}
			}
		}
	default:
		log.Println("usage: qobuz-sync <album|track|favorites>")

		return
	}
}

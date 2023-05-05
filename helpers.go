package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/frolovo22/tag"
	"github.com/pkg/errors"
)

//nolint:gochecknoglobals
var (
	albumTracker *Track
	trackTracker *Track

	dirPerm  = os.FileMode(0o755) //nolint:gomnd
	filePerm = os.FileMode(0o644) //nolint:gomnd

	infoMsg = "\x1b[34minfo:\x1b[0m"
	warnMsg = "\x1b[33mwarn:\x1b[0m"
	errMsg  = "\x1b[31merr:\x1b[0m"
)

type tags struct {
	Album       string
	AlbumArtist string
	Artist      string
	Comment     string
	Composer    string
	Genre       string
	Date        time.Time
	Disc        int
	DiscTotal   int
	Track       int
	TrackTotal  int
	Title       string
}

func sanitizeStringToPath(str string) string {
	str = strings.ReplaceAll(str, "/", "_")
	str = strings.ReplaceAll(str, "\\", "_")
	str = strings.ReplaceAll(str, ":", "_")
	str = strings.ReplaceAll(str, "*", "_")
	str = strings.ReplaceAll(str, "?", "_")
	str = strings.ReplaceAll(str, "\"", "_")
	str = strings.ReplaceAll(str, "<", "_")
	str = strings.ReplaceAll(str, ">", "_")
	str = strings.ReplaceAll(str, "|", "_")

	str = strings.TrimSpace(str)

	return str
}

func buildPath(dir string, title string, track int, trackTotal int, disc int, discTotal int) string {
	title = sanitizeStringToPath(title)
	padding := len(strconv.Itoa(trackTotal))
	trackNumber := fmt.Sprintf("%0"+strconv.Itoa(padding)+"d", track)
	discNumber := ""

	if discTotal > 1 {
		discNumber = fmt.Sprintf("%d", disc)
	}

	format := ".flac"

	return filepath.Join(dir, fmt.Sprintf("%v%v - %v%v", discNumber, trackNumber, title, format))
}

func setTags(path string, t tags) error {
	tags, err := tag.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read part file for tags")
	}

	for _, err := range []error{
		tags.SetAlbum(t.Album),
		tags.SetAlbumArtist(t.AlbumArtist),
		tags.SetArtist(t.Artist),

		tags.SetComment(t.Comment),
		tags.SetGenre(t.Genre),
		tags.SetDate(t.Date),
		tags.SetComposer(t.Composer),

		tags.SetTitle(t.Title),
		tags.SetTrackNumber(t.Track, t.TrackTotal),
		tags.SetDiscNumber(t.Disc, t.DiscTotal),
	} {
		if err != nil {
			return errors.Wrap(err, "failed to set tag")
		}
	}

	err = tags.SaveFile(strings.ReplaceAll(path, ".part", ""))
	if err != nil {
		return errors.Wrap(err, "failed to save tags")
	}

	err = os.Remove(path)
	if err != nil {
		return errors.Wrap(err, "failed to remove part file")
	}

	return nil
}

func downloadTrackToFile(c *Client, trackID string, path string) error {
	url, err := c.TrackGetFileURL(trackID, QualityMAX)
	if err != nil {
		return errors.Wrap(err, "failed to get track file url")
	}

	if url.FormatID == 0 {
		return ErrUnavailable
	}

	if url.MimeType != "audio/flac" {
		log.Printf(warnMsg, "not FLAC got %v: %v\n", url.MimeType, path)
		path = strings.ReplaceAll(path, ".flac", ".mp3")
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

	f, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
}

func downloadTrack(c *Client, id string, dir string) (string, error) {
	path, err := trackTracker.Get(id)
	if err == nil {
		return path, errors.Wrap(ErrAlreadyExists, "cached")
	}

	res, err := c.TrackGet(id)
	if err != nil {
		return "", errors.Wrap(err, "failed to get track")
	}

	if dir == "" {
		return "", errors.New("`dir` argument not set")
	}

	err = os.MkdirAll(dir, dirPerm)
	if err != nil {
		return "", errors.Wrap(err, "failed to create dir")
	}

	url, err := c.TrackGetFileURL(id, QualityMAX)
	if err != nil {
		return "", errors.Wrap(err, "failed to get track file url")
	}

	if url.FormatID == 0 {
		return "", ErrUnavailable
	}

	path = buildPath(dir, res.Title, res.TrackNumber, res.Album.TracksCount, res.MediaNumber, res.Album.MediaCount)

	err = os.MkdirAll(filepath.Dir(path), dirPerm)
	if err != nil {
		return "", errors.Wrap(err, "failed to create dir")
	}

	_, err = os.Stat(path)
	if err == nil {
		err = trackTracker.Set(id, path)
		if err != nil {
			return "", errors.Wrap(err, "failed to set track as downloaded")
		}

		return path, ErrAlreadyExists
	}

	p := path + ".part"

	err = downloadTrackToFile(c, strconv.Itoa(res.ID), p)
	if err != nil {
		return "", errors.Wrap(err, "failed to download track")
	}

	err = setTags(p, tags{
		Album:       res.Album.Title,
		AlbumArtist: res.Album.Artist.Name,
		Artist:      res.Performer.Name,
		Comment:     fmt.Sprintf("qobuz_id: %v", res.ID),
		Composer:    res.Composer.Name,
		Genre:       res.Album.Genre.Name,
		Date:        time.Unix(int64(res.Album.ReleasedAt), 0),
		Disc:        res.MediaNumber,
		DiscTotal:   res.Album.MediaCount,
		Track:       res.TrackNumber,
		TrackTotal:  res.Album.TracksCount,
		Title:       res.Title,
	})
	if err != nil {
		return "", errors.Wrap(err, "failed to set tags")
	}

	log.Println(infoMsg, "downloaded track", path)

	err = trackTracker.Set(id, path)
	if err != nil {
		return "", errors.Wrap(err, "failed to set track as downloaded")
	}

	return path, nil
}

func downloadAlbumArt(url string, dir string) error {
	_, err := os.Stat(filepath.Join(dir, "album.jpg"))
	if err == nil {
		return ErrAlreadyExists
	}

	u := strings.ReplaceAll(url, "_600", "_org")
	if u != url {
		err := downloadAlbumArt(u, dir)
		if err == nil {
			return nil
		}
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get album art")
	}

	defer res.Body.Close()

	ff, err := os.Create(filepath.Join(dir, "album.jpg"))
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	_, err = io.Copy(ff, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
}

func downloadAlbum(c *Client, id string, dir string) (string, error) {
	d, err := albumTracker.Get(id)
	if err == nil {
		return d, errors.Wrap(ErrAlreadyExists, "cached")
	}

	res, err := c.AlbumGet(id)
	artist := sanitizeStringToPath(res.Artist.Name)
	albumName := sanitizeStringToPath(res.Title)
	dir = filepath.Join(dir, artist, albumName)

	if err != nil {
		return dir, errors.Wrap(err, "album get")
	}

	for _, track := range res.Tracks.Items {
		path := buildPath(dir, track.Title, track.TrackNumber, res.TracksCount, track.MediaNumber, res.MediaCount)

		_, err = os.Stat(path)
		if err == nil {
			log.Println(infoMsg, "track already exists, skipping:", path)

			continue
		}

		_, err := downloadTrack(c, strconv.Itoa(track.ID), dir)
		if err != nil {
			if errors.Is(err, ErrAlreadyExists) {
				log.Println(infoMsg, "track already exists, skipping:", path)
			} else {
				log.Println(warnMsg, "failed to download track, skipping:", err)
			}

			continue
		}
	}

	err = downloadAlbumArt(res.Image.Large, dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Println(infoMsg, "album art already exists, skipping:", dir)
		} else {
			log.Println(warnMsg, "failed to download album art, skipping:", err)
		}
	}

	log.Println(infoMsg, "downloaded album", dir)

	err = albumTracker.Set(id, dir)
	if err != nil {
		return dir, errors.Wrap(err, "failed to set album as downloaded")
	}

	return dir, nil
}

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
	AlbumArt    string
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
		log.Printf("warn: not FLAC got %v: %v\n", url.MimeType, path)
		path = strings.ReplaceAll(path, ".flac", ".mp3")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url.URL, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	client := &http.Client{}
	res, err := client.Do(req)
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
	res, err := c.TrackGet(id)
	if err != nil {
		return "", errors.Wrap(err, "failed to get track")
	}

	if dir == "" {
		return "", errors.New("`dir` argument not set")
	}

	err = os.MkdirAll(dir, 0o755)
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

	path := buildPath(dir, res.Title, res.TrackNumber, res.Album.TracksCount, res.MediaNumber, res.Album.MediaCount)
	err = os.MkdirAll(filepath.Dir(path), 0o755)
	if err != nil {
		return "", errors.Wrap(err, "failed to create dir")
	}

	_, err = os.Stat(path)
	if err == nil {
		return path, ErrAlreadyExists
	}

	p := path + ".part"
	err = downloadTrackToFile(c, strconv.Itoa(res.ID), p)
	if err != nil {
		return "", errors.Wrap(err, "failed to download track")
	}

	desc := fmt.Sprintf("qobuz_id: %v", res.ID)
	err = setTags(p, tags{
		Album:       res.Album.Title,
		AlbumArtist: res.Album.Artist.Name,
		Artist:      res.Performer.Name,
		Comment:     desc,
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

	log.Println("info: downloaded track", path)

	return path, nil
}

func downloadAlbumArt(url string, dir string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	client := &http.Client{}
	res, err := client.Do(req)
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

func downloadAlbum(c *Client, id string, dir string) error {
	res, err := c.AlbumGet(id)
	if err != nil {
		return errors.Wrap(err, "album get")
	}

	artist := sanitizeStringToPath(res.Artist.Name)
	albumName := sanitizeStringToPath(res.Title)
	dir = filepath.Join(dir, artist, albumName)

	for _, track := range res.Tracks.Items {
		path := buildPath(dir, track.Title, track.TrackNumber, res.TracksCount, track.MediaNumber, res.MediaCount)
		_, err = os.Stat(path)
		if err == nil {
			log.Println("info: track already exists, skipping:", path)

			continue
		}

		_, err := downloadTrack(c, strconv.Itoa(track.ID), dir)
		if err != nil {
			if errors.Is(err, ErrAlreadyExists) {
				log.Println("info: track already exists, skipping:", path)
			} else {
				log.Println("warn: failed to download track, skipping:", err)
			}

			continue
		}
	}

	err = downloadAlbumArt(res.Image.Large, dir)
	if err != nil {
		log.Println("warn: failed to download album art, skipping:", err)
	}

	log.Println("info: downloaded album", dir)

	return nil
}

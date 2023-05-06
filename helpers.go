package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/frolovo22/tag"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

//nolint:gochecknoglobals
var (
	dirPerm  = os.FileMode(0o755) //nolint:gomnd
	filePerm = os.FileMode(0o644) //nolint:gomnd

	albumTracker *Tracker
	trackTracker *Tracker
)

type MetadataTags struct {
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

func setTags(path string, tags MetadataTags) error {
	fileTags, err := tag.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "failed to read part file for tags")
	}

	for _, err := range []error{
		fileTags.SetAlbum(tags.Album),
		fileTags.SetAlbumArtist(tags.AlbumArtist),
		fileTags.SetArtist(tags.Artist),

		fileTags.SetComment(tags.Comment),
		fileTags.SetGenre(tags.Genre),
		fileTags.SetDate(tags.Date),
		fileTags.SetComposer(tags.Composer),

		fileTags.SetTitle(tags.Title),
		fileTags.SetTrackNumber(tags.Track, tags.TrackTotal),
		fileTags.SetDiscNumber(tags.Disc, tags.DiscTotal),
	} {
		if err != nil {
			return errors.Wrap(err, "failed to set tag")
		}
	}

	err = fileTags.SaveFile(strings.ReplaceAll(path, ".part", ""))
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
		log.Warn().Msgf("not FLAC got %v: %v", url.MimeType, path)
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

//nolint:cyclop,funlen // todo: fix in future
func downloadTrack(client *Client, trackID string, dir string) (string, error) {
	path, err := trackTracker.Get(trackID)
	if err == nil {
		return path, errors.Wrap(ErrAlreadyExists, "cached")
	}

	res, err := client.TrackGet(trackID)
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

	path = buildPath(dir, res.Title, res.TrackNumber, res.Album.TracksCount, res.MediaNumber, res.Album.MediaCount)

	err = os.MkdirAll(filepath.Dir(path), dirPerm)
	if err != nil {
		return "", errors.Wrap(err, "failed to create dir")
	}

	_, err = os.Stat(path)
	if err == nil {
		err = trackTracker.Set(trackID, path)
		if err != nil {
			return "", errors.Wrap(err, "failed to set track as downloaded")
		}

		return path, ErrAlreadyExists
	}

	partialPath := path + ".part"

	err = downloadTrackToFile(client, strconv.Itoa(res.ID), partialPath)
	if err != nil {
		return "", errors.Wrap(err, "failed to download track")
	}

	err = setTags(partialPath, MetadataTags{
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

	log.Info().Msgf("downloaded track: %v", path)

	err = trackTracker.Set(trackID, path)
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

	if strings.Contains(url, "_600") {
		err := downloadAlbumArt(strings.ReplaceAll(url, "_600", "_org"), dir)
		if err == nil {
			return nil
		}
	}

	log.Debug().Str("url", url).Str("dir", dir).Msg("downloading album art")

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

func downloadAlbum(client *Client, albumID string, dir string) (string, error) {
	d, err := albumTracker.Get(albumID)
	if err == nil {
		return d, errors.Wrap(ErrAlreadyExists, "cached")
	}

	res, err := client.AlbumGet(albumID)
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
			log.Info().Msgf("track already exists, skipping: %v", path)

			continue
		}

		_, err := downloadTrack(client, strconv.Itoa(track.ID), dir)
		if err != nil {
			if errors.Is(err, ErrAlreadyExists) {
				log.Info().Msgf("track already exists, skipping: %v", path)
			} else {
				log.Warn().Msgf("failed to download track, skipping: %v", err)
			}

			continue
		}
	}

	err = downloadAlbumArt(res.Image.Large, dir)
	if err != nil {
		if errors.Is(err, ErrAlreadyExists) {
			log.Info().Msgf("album art already exists, skipping: %v", dir)
		} else {
			log.Warn().Msgf("failed to download album art, skipping: %v", err)
		}
	}

	log.Info().Msgf("downloaded album: %v", dir)

	err = albumTracker.Set(albumID, dir)
	if err != nil {
		return dir, errors.Wrap(err, "failed to set album as downloaded")
	}

	return dir, nil
}

func getHandlerFromContext(ctx context.Context) (*Handler, error) {
	//nolint:exhaustruct
	switch t := ctx.Value(Handler{}).(type) {
	case *Handler:
		return t, nil
	case nil:
		return nil, errors.New("handler is nil")
	default:
		return nil, errors.New("handler is not a *Handler")
	}
}

func envOrDefault(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

package responses

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/helpers"
)

type Track struct {
	MaximumBitDepth int    `json:"maximum_bit_depth"`
	Copyright       string `json:"copyright"`
	Performers      string `json:"performers"`
	AudioInfo       struct {
		ReplaygainTrackPeak float64 `json:"replaygain_track_peak"`
		ReplaygainTrackGain float64 `json:"replaygain_track_gain"`
	} `json:"audio_info"`
	Performer           *Artist `json:"performer"`
	Album               *Album  `json:"album"`
	Work                any     `json:"work"`
	Composer            *Artist `json:"composer"`
	Isrc                string  `json:"isrc"`
	Title               string  `json:"title"`
	Version             string  `json:"version"`
	Duration            int     `json:"duration"`
	ParentalWarning     bool    `json:"parental_warning"`
	TrackNumber         int     `json:"track_number"`
	MaximumChannelCount int     `json:"maximum_channel_count"`
	ID                  int     `json:"id"`
	MediaNumber         int     `json:"media_number"`
	MaximumSamplingRate float64 `json:"maximum_sampling_rate"`
	ReleaseDateOriginal any     `json:"release_date_original"`
	ReleaseDateDownload any     `json:"release_date_download"`
	ReleaseDateStream   any     `json:"release_date_stream"`
	Purchasable         bool    `json:"purchasable"`
	Streamable          bool    `json:"streamable"`
	Previewable         bool    `json:"previewable"`
	Sampleable          bool    `json:"sampleable"`
	Downloadable        bool    `json:"downloadable"`
	Displayable         bool    `json:"displayable"`
	PurchasableAt       int     `json:"purchasable_at"`
	StreamableAt        int     `json:"streamable_at"`
	Hires               bool    `json:"hires"`
	HiresStreamable     bool    `json:"hires_streamable"`
	Position            int     `json:"position"`
	CreatedAt           int     `json:"created_at"`
	PlaylistTrackID     int     `json:"playlist_track_id"`
	FavoritedAt         int     `json:"favorited_at"`
}

func (t Track) Filename() string {
	title := helpers.SanitizeStringToPath(t.Title)
	padding := len(strconv.Itoa(t.Album.TracksCount))
	trackNumber := fmt.Sprintf("%0"+strconv.Itoa(padding)+"d", t.TrackNumber)
	discNumber := ""

	if t.Album.MediaCount > 1 {
		discNumber = strconv.Itoa(t.MediaNumber)
	}

	format := ".flac"

	return fmt.Sprintf("%v%v - %v%v", discNumber, trackNumber, title, format)
}

func (t Track) Path() string {
	if t.Album == nil {
		spew.Dump(t)
	}

	path := filepath.Join(t.Album.Path(), t.Filename())
	path = filepath.Clean(path)

	return path
}

func (t Track) Metadata() common.Metadata {
	if t.Album == nil {
		return common.Metadata{}
	}

	metadata := common.Metadata{
		Album:       t.Album.Title,
		AlbumArtist: t.Album.Artist.Name,
		Comment:     fmt.Sprintf("qobuz_id: %v", t.ID),
		Genre:       t.Album.Genre.Name,
		Date:        time.Unix(int64(t.Album.ReleasedAt), 0),
		Disc:        t.MediaNumber,
		DiscTotal:   t.Album.MediaCount,
		Track:       t.TrackNumber,
		TrackTotal:  t.Album.TracksCount,
		Title:       t.Title,
	}

	if t.Composer != nil {
		metadata.Composer = t.Composer.Name
	} else {
		metadata.Composer = ""
	}

	if t.Performer != nil {
		metadata.Artist = t.Performer.Name
	} else if t.Performers != "" {
		metadata.Artist = t.Performers
	} else if t.Album.Artist != nil {
		metadata.Artist = t.Album.Artist.Name
	}

	return metadata
}

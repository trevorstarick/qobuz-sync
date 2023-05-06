package main

import (
	"crypto/md5" //nolint:gosec // MD5 is used for request signatures, not security
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"

	albumGet "github.com/trevorstarick/qobuz-sync/responses/album/get"
	favoriteGetUserFavorites "github.com/trevorstarick/qobuz-sync/responses/favorite/getUserFavorites"
	trackGet "github.com/trevorstarick/qobuz-sync/responses/track/get"
	trackGetFileUrl "github.com/trevorstarick/qobuz-sync/responses/track/getFileUrl"
	trackSearch "github.com/trevorstarick/qobuz-sync/responses/track/search"
)

type trackFormat int

const (
	QualityMP3   trackFormat = 5  // 320kbps
	QualityFLAC  trackFormat = 6  // 16-bit 44.1kHz+
	QualityHIRES trackFormat = 7  // 24-bit 44.1kHz+
	QualityMAX   trackFormat = 27 // 24-bit 96kHz+

)

type listType string

const (
	ListTypeALBUM  listType = "albums"
	ListTypeTRACK  listType = "tracks"
	listTypeARTIST listType = "artists"
)

func (c *Client) TrackSearch(query string) (*trackSearch.TrackSearch, error) {
	return (Querier[trackSearch.TrackSearch]{c}).Req("track/search", &url.Values{
		"query": []string{query},
	})
}

func (c *Client) TrackGetFileURL(trackID string, format trackFormat) (*trackGetFileUrl.Response, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sig := "trackgetFileUrlformat_id%vintentstreamtrack_id%v%v%v"
	sig = fmt.Sprintf(sig, format, trackID, timestamp, c.Secrets[0])
	hash := md5.Sum([]byte(sig)) //nolint:gosec // MD5 is used for request signatures, not security
	hashedSig := hex.EncodeToString(hash[:])

	return (Querier[trackGetFileUrl.Response]{c}).Req("track/getFileUrl", &url.Values{
		"request_ts":  []string{timestamp},
		"request_sig": []string{hashedSig},
		"track_id":    []string{trackID},
		"format_id":   []string{strconv.Itoa(int(format))},
		"intent":      []string{"stream"},
	})
}

func (c *Client) TrackGet(trackID string) (*trackGet.Response, error) {
	return (Querier[trackGet.Response]{c}).Req("track/get", &url.Values{
		"track_id": []string{trackID},
	})
}

func (c *Client) FavoriteGetUserFavorites(listType listType, offset int) (*favoriteGetUserFavorites.Response, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sig := "favoritegetUserFavorites" + timestamp
	hash := md5.Sum([]byte(sig)) //nolint:gosec // MD5 is used for request signatures, not security
	hashedSig := hex.EncodeToString(hash[:])

	return (Querier[favoriteGetUserFavorites.Response]{c}).Req("favorite/getUserFavorites", &url.Values{
		"limit":       []string{"100"},
		"offset":      []string{strconv.Itoa(offset)},
		"type":        []string{string(listType)}, // albums, tracks, artists, article
		"request_ts":  []string{timestamp},
		"request_sig": []string{hashedSig},
	})
}

func (c *Client) AlbumGet(albumID string) (*albumGet.Response, error) {
	return (Querier[albumGet.Response]{c}).Req("album/get", &url.Values{
		"album_id": []string{albumID},
	})
}

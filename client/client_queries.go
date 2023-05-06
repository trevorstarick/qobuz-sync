package client

import (
	"crypto/md5" //nolint:gosec // MD5 is used for request signatures
	"encoding/hex"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/responses"
	albumGet "github.com/trevorstarick/qobuz-sync/responses/album/get"
	favoriteGetUserFavorites "github.com/trevorstarick/qobuz-sync/responses/favorite/getUserFavorites"
	playlistGet "github.com/trevorstarick/qobuz-sync/responses/playlist/get"
	trackGet "github.com/trevorstarick/qobuz-sync/responses/track/get"
	trackGetFileUrl "github.com/trevorstarick/qobuz-sync/responses/track/getFileUrl"
	trackSearch "github.com/trevorstarick/qobuz-sync/responses/track/search"
)

func (client *Client) TrackSearch(query string) (*trackSearch.TrackSearch, error) {
	return (Querier[trackSearch.TrackSearch]{client}).Req("track/search", &url.Values{
		"query": []string{query},
	})
}

func (client *Client) TrackGetFileURL(trackID string, format trackFormat) (*trackGetFileUrl.Response, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sig := "trackgetFileUrlformat_id%vintentstreamtrack_id%v%v%v"
	sig = fmt.Sprintf(sig, format, trackID, timestamp, client.Secrets[0])
	hash := md5.Sum([]byte(sig)) //nolint:gosec // MD5 is used for request signatures, not security
	hashedSig := hex.EncodeToString(hash[:])

	res, err := (Querier[trackGetFileUrl.Response]{client}).Req("track/getFileUrl", &url.Values{
		"request_ts":  []string{timestamp},
		"request_sig": []string{hashedSig},
		"track_id":    []string{trackID},
		"format_id":   []string{strconv.Itoa(int(format))},
		"intent":      []string{"stream"},
	})
	if err != nil {
		return nil, err
	}

	if res.FormatID == 0 {
		return nil, common.ErrUnavailable
	}

	return res, nil
}

func (client *Client) TrackGet(trackID string) (*trackGet.Response, error) {
	return (Querier[trackGet.Response]{client}).Req("track/get", &url.Values{
		"track_id": []string{trackID},
	})
}

func (client *Client) FavoriteGetUserFavorites(listType listType, offset int) (*favoriteGetUserFavorites.Response, error) { //nolint:lll // long function name
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	sig := "favoritegetUserFavorites" + timestamp
	hash := md5.Sum([]byte(sig)) //nolint:gosec // MD5 is used for request signatures, not security
	hashedSig := hex.EncodeToString(hash[:])

	return (Querier[favoriteGetUserFavorites.Response]{client}).Req("favorite/getUserFavorites", &url.Values{
		"limit":       []string{"100"},
		"offset":      []string{strconv.Itoa(offset)},
		"type":        []string{string(listType)}, // albums, tracks, artists, article
		"request_ts":  []string{timestamp},
		"request_sig": []string{hashedSig},
	})
}

func (client *Client) AlbumGet(albumID string) (*albumGet.Response, error) {
	return (Querier[albumGet.Response]{client}).Req("album/get", &url.Values{
		"album_id": []string{albumID},
	})
}

func (client *Client) PlaylistGet(playlistID string) (*playlistGet.Response, error) {
	var (
		res *playlistGet.Response
		err error

		limit  = 500
		offset = 0
		tracks = make([]responses.Track, 0)
	)

	for {
		res, err = (Querier[playlistGet.Response]{client}).Req("playlist/get", &url.Values{
			"playlist_id": []string{playlistID},
			"extra":       []string{"tracks"},
			"limit":       []string{strconv.Itoa(limit)},
			"offset":      []string{strconv.Itoa(offset)},
		})
		if err != nil {
			return nil, err
		}

		tracks = append(tracks, res.Tracks.Items...)

		if res.Tracks.Total < res.Tracks.Offset+res.Tracks.Limit {
			break
		}
	}

	res.Tracks.Items = tracks

	return res, nil
}

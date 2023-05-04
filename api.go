package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
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
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"track/search", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	q := req.URL.Query()
	q.Set("query", query)
	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var response trackSearch.TrackSearch

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &response, nil
}

func (c *Client) TrackGetFileURL(trackID string, format trackFormat) (*trackGetFileUrl.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"track/getFileUrl", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())

	sig := "trackgetFileUrlformat_id%vintentstreamtrack_id%v%v%v"
	sig = fmt.Sprintf(sig, format, trackID, ts, c.Secrets[0])

	hash := md5.Sum([]byte(sig))
	hashedSig := hex.EncodeToString(hash[:])

	q := req.URL.Query()
	q.Set("request_ts", ts)
	q.Set("request_sig", hashedSig)
	q.Set("track_id", trackID)
	q.Set("format_id", fmt.Sprintf("%d", format))
	q.Set("intent", "stream")

	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var response trackGetFileUrl.Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &response, nil
}

func (c *Client) TrackGet(trackID string) (*trackGet.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"track/get", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	q := req.URL.Query()
	q.Set("track_id", trackID)

	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var response trackGet.Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &response, nil
}

func (c *Client) FavoriteGetUserFavorites(t listType, offset int) (*favoriteGetUserFavorites.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"favorite/getUserFavorites", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	ts := fmt.Sprintf("%d", time.Now().Unix())
	sig := "favoritegetUserFavorites" + ts
	hash := md5.Sum([]byte(sig))
	hashedSig := hex.EncodeToString(hash[:])

	q := req.URL.Query()
	q.Set("limit", "100")
	q.Set("offset", strconv.Itoa(offset))
	q.Set("type", string(t)) // albums, tracks, artists, articles
	q.Set("request_ts", ts)
	q.Set("request_sig", hashedSig)

	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		buf, err := io.ReadAll(res.Body)
		_ = err
		fmt.Fprintf(os.Stderr, "req: %v\n", req.URL.String())
		fmt.Fprintf(os.Stderr, "err: %v\n", string(buf))

		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var response favoriteGetUserFavorites.Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &response, nil
}

func (c *Client) AlbumGet(albumID string) (*albumGet.Response, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"album/get", nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	q := req.URL.Query()
	q.Set("album_id", albumID)

	req.URL.RawQuery = q.Encode()

	res, err := c.do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var response albumGet.Response

	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, errors.Wrap(err, "failed to decode response")
	}

	return &response, nil
}

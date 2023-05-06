package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
)

type Querier[T any] struct {
	*Client
}

func (q Querier[T]) prepareRequest(path string, query *url.Values) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+path, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create request")
	}

	qq := req.URL.Query()
	for k, v := range *query {
		qq.Set(k, v[0])
	}

	req.URL.RawQuery = qq.Encode()

	for k, v := range q.Header {
		req.Header.Set(k, v[0])
	}

	queryParams := req.URL.Query()
	queryParams.Set("app_id", q.AppID)
	req.URL.RawQuery = queryParams.Encode()

	return req, nil
}

func (q Querier[T]) Req(path string, query *url.Values) (*T, error) {
	log.Debug().Str("path", path).Str("query", query.Encode()).Msg("requesting")

	req, err := q.prepareRequest(path, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare request")
	}

	res, err := q.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = errors.Wrap(err, "failed to close m3u file")
		}
	}()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, errors.Wrap(common.ErrNotFound, "not found")
	case http.StatusBadRequest:
		return nil, errors.Wrap(common.ErrBadRequest, "bad request")
	case http.StatusUnauthorized:
		return nil, errors.Wrap(common.ErrAuthFailed, "invalid token")
	default:
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	t := new(T)
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return t, nil
}

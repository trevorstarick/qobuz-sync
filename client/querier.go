package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Querier[T any] struct {
	*Client
}

func (q Querier[T]) Req(path string, query *url.Values) (*T, error) {
	log.Debug().Str("path", path).Str("query", query.Encode()).Msg("requesting")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+path, nil)
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

	res, err := q.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	t := new(T)
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		return nil, errors.Wrap(err, "failed to decode response body")
	}

	return t, nil
}

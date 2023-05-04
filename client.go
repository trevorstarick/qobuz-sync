package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	userLogin "github.com/trevorstarick/qobuz-sync/responses/user/login"
)

const (
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0"
	baseApp       = "https://play.qobuz.com"
	baseAPI       = "https://www.qobuz.com/api.json/0.2/"
	userAuthToken = "X-User-Auth-Token"
)

var (
	bundleRegexp = regexp.MustCompile(`(/resources/\d+\.\d+\.\d+-[a-z]\d{3}/bundle\.js)`)
	appIDRegexp  = regexp.MustCompile(`production:{api:{appId:"(?P<app_id>\d{9})",appSecret:"\w{32}"`)
	tzRegexp     = regexp.MustCompile(`[a-z]\.initialSeed\("(?P<seed>[\w=]+)",window\.utimezone\.(?P<timezone>[a-z]+)\)`)
	infoRegexp   = regexp.MustCompile(`name:"\w+/([A-Z][a-z]+)",info:"(?P<info>[\w=]+)",extras:"(?P<extras>[\w=]+)"`)
)

type Client struct {
	c *http.Client `json:"-"`

	AppID   string
	Secrets []string
	Header  http.Header
}

func (c *Client) getBundleURL() (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseApp+"/login", nil)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	matches := bundleRegexp.FindAllStringSubmatch(string(buf), -1)
	if len(matches) == 0 {
		return "", errors.New("no matches found")
	}

	return baseApp + matches[0][1], nil
}

func (c *Client) getAppID() (string, error) {
	bundleURL, err := c.getBundleURL()
	if err != nil {
		return "", errors.Wrap(err, "get bundle url")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, bundleURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	matches := appIDRegexp.FindAllStringSubmatch(string(buf), -1)
	if len(matches) == 0 {
		return "", errors.New("no matches found")
	}

	return matches[0][1], nil
}

func (c *Client) testSecret(secret string) bool {
	secrets := c.Secrets
	c.Secrets = []string{secret}
	defer func() {
		c.Secrets = secrets
	}()

	_, err := c.TrackGetFileURL("5966783", QualityMP3)

	return err == nil
}

func (c *Client) getSecrets() ([]string, error) {
	bundleURL, err := c.getBundleURL()
	if err != nil {
		return nil, errors.Wrap(err, "get bundle url")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, bundleURL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "new request")
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "do request")
	}

	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.Wrap(err, "read body")
	}

	matches := tzRegexp.FindAllStringSubmatch(string(buf), -1)
	if len(matches) == 0 {
		return nil, errors.New("no matches found")
	}

	type t struct {
		seed   string
		info   string
		extras string
	}

	seeds := make(map[string]t)
	for _, match := range matches {
		seeds[match[2]] = t{
			seed: match[1],
		}
	}

	infos := infoRegexp.FindAllStringSubmatch(string(buf), -1)
	if len(infos) == 0 {
		return nil, errors.New("no matches found")
	}

	for _, info := range infos {
		tz := strings.ToLower(info[1])
		if tz == "algiers" {
			tz = "algier"
		}

		if _, ok := seeds[tz]; ok {
			seed := seeds[tz].seed
			seeds[tz] = t{
				seed:   seed,
				info:   info[2],
				extras: info[3],
			}
		}
	}

	secrets := make([]string, 0)

	for _, seed := range seeds {
		rightSide := seed.seed + seed.info + seed.extras
		secret := rightSide[:len(rightSide)-44]
		base64Secret, err := base64.StdEncoding.DecodeString(secret)
		if err != nil {
			return nil, errors.Wrap(err, "decode secret")
		}

		if c.testSecret(string(base64Secret)) {
			secrets = append(secrets, string(base64Secret))
		}
	}

	return secrets, nil
}

func NewClient(email, password string) (*Client, error) {
	headers := http.Header{}
	headers.Set("User-Agent", userAgent)

	client := &Client{
		c: &http.Client{},
	}

	appID, err := client.getAppID()
	if err != nil {
		return nil, errors.Wrap(err, "get app id")
	}

	if appID == "" {
		return nil, errors.New("no app id found")
	}

	client = &Client{
		c:      &http.Client{},
		AppID:  appID,
		Header: headers,
	}

	if err := client.auth(email, password); err != nil {
		return nil, errors.Wrap(err, "auth")
	}

	secrets, err := client.getSecrets()
	if err != nil {
		return nil, errors.Wrap(err, "get secrets")
	}

	if secrets == nil {
		return nil, errors.New("no secrets found")
	}

	client.Secrets = secrets

	return client, nil
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	for k, v := range c.Header {
		req.Header.Set(k, v[0])
	}

	q := req.URL.Query()
	q.Set("app_id", c.AppID)
	req.URL.RawQuery = q.Encode()

	res, err := c.c.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do request")
	}

	return res, nil
}

func (c *Client) auth(email, password string) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseAPI+"user/login", nil)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	query := req.URL.Query()
	query.Set("email", email)
	query.Set("password", password)
	query.Set("app_id", c.AppID)
	req.URL.RawQuery = query.Encode()

	res, err := c.do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do request")
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status code: %d", res.StatusCode)
	}

	var login userLogin.Response

	if err := json.NewDecoder(res.Body).Decode(&login); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if login.UserAuthToken == "" {
		return errors.New("no user auth token found")
	}

	c.Header.Set(userAuthToken, login.UserAuthToken)

	return nil
}

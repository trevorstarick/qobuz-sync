package main

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	userLogin "github.com/trevorstarick/qobuz-sync/responses/user/login"
)

const (
	userAgent     = "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:83.0) Gecko/20100101 Firefox/83.0"
	baseApp       = "https://play.qobuz.com"
	baseAPI       = "https://www.qobuz.com/api.json/0.2/"
	userAuthToken = "X-User-Auth-Token" //nolint:gosec // This is not a secret
)

var (
	bundleRegexp = regexp.MustCompile(`(/resources/\d+\.\d+\.\d+-[a-z]\d{3}/bundle\.js)`)
	appIDRegexp  = regexp.MustCompile(`production:{api:{appId:"(?P<app_id>\d{9})",appSecret:"\w{32}"`)
	tzRegexp     = regexp.MustCompile(`[a-z]\.initialSeed\("(?P<seed>[\w=]+)",window\.utimezone\.(?P<timezone>[a-z]+)\)`)
	infoRegexp   = regexp.MustCompile(`name:"\w+/([A-Z][a-z]+)",info:"(?P<info>[\w=]+)",extras:"(?P<extras>[\w=]+)"`)
)

type Client struct {
	c *http.Client `json:"-"`

	bundle string

	AppID   string
	Secrets []string
	Header  http.Header
}

func (c *Client) getBundle() (string, error) {
	if c.bundle != "" {
		return c.bundle, nil
	}

	bundleURL, err := c.getBundleURL()
	if err != nil {
		return "", errors.Wrap(err, "get bundle url")
	}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, bundleURL, nil)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	c.bundle = string(buf)

	return c.bundle, nil
}

func (c *Client) getBundleURL() (string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, baseApp+"/login", nil)
	if err != nil {
		return "", errors.Wrap(err, "new request")
	}

	res, err := http.DefaultClient.Do(req)
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
	bundle, err := c.getBundle()
	if err != nil {
		return "", errors.Wrap(err, "get bundle")
	}

	matches := appIDRegexp.FindAllStringSubmatch(bundle, -1)
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

//nolint:cyclop,funlen // todo: fix in future
func (c *Client) getSecrets() ([]string, error) {
	bundle, err := c.getBundle()
	if err != nil {
		return nil, errors.Wrap(err, "get bundle")
	}

	matches := tzRegexp.FindAllStringSubmatch(bundle, -1)
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
			seed:   match[1],
			info:   "",
			extras: "",
		}
	}

	infos := infoRegexp.FindAllStringSubmatch(bundle, -1)
	if len(infos) == 0 {
		return nil, errors.New("no matches found")
	}

	for _, info := range infos {
		timezone := strings.ToLower(info[1])
		if timezone == "algiers" {
			timezone = "algier"
		}

		if _, ok := seeds[timezone]; ok {
			seed := seeds[timezone].seed
			seeds[timezone] = t{
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
		c:       http.DefaultClient,
		bundle:  "",
		AppID:   "",
		Header:  headers,
		Secrets: []string{},
	}

	appID, err := client.getAppID()
	if err != nil {
		return nil, errors.Wrap(err, "get app id")
	}

	if appID == "" {
		return nil, errors.New("no app id found")
	}

	client.AppID = appID

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

func (c *Client) auth(email, password string) error {
	login, err := (Querier[userLogin.Response]{c}).Req("user/login", &url.Values{
		"email":    {email},
		"password": {password},
		"app_id":   {c.AppID},
	})
	if err != nil {
		return errors.Wrap(err, "user login")
	}

	if login.UserAuthToken == "" {
		return errors.New("no user auth token found")
	}

	c.Header.Set(userAuthToken, login.UserAuthToken)

	return nil
}

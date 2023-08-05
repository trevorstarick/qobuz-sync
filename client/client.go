package client

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
	userLogin "github.com/trevorstarick/qobuz-sync/responses/user/login"
)

type trackFormat int

const (
	QualityMP3   trackFormat = 5  // 320kbps
	QualityFLAC  trackFormat = 6  // 16-bit 44.1kHz+
	QualityHIRES trackFormat = 7  // 24-bit 44.1kHz+
	QualityMAX   trackFormat = 27 // 24-bit 96kHz+

)

type ListType string

const (
	ListTypeALBUM  ListType = "albums"
	ListTypeTRACK  ListType = "tracks"
	listTypeARTIST ListType = "artists"
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

type Key struct{}

type Client struct {
	c *http.Client

	baseDir string

	albumTracker *Tracker
	trackTracker *Tracker

	bundle string

	force bool

	AppID   string
	Secrets []string
	Header  http.Header
}

func NewClient(email, password, baseDir string, force bool) (*Client, error) {
	headers := http.Header{}
	headers.Set("User-Agent", userAgent)

	client := &Client{
		c:            http.DefaultClient,
		bundle:       "",
		baseDir:      baseDir,
		trackTracker: &Tracker{}, //nolint:exhaustruct
		albumTracker: &Tracker{}, //nolint:exhaustruct
		force:        force,
		AppID:        "",
		Header:       headers,
		Secrets:      []string{},
	}

	appID, err := client.getAppID()
	if err != nil {
		return nil, errors.Wrap(err, "get app id")
	}

	if appID == "" {
		return nil, errors.New("no app id found")
	}

	client.AppID = appID

	if err := client.Login(email, password); err != nil {
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

	client.trackTracker, err = NewTracker(filepath.Join(baseDir, "tracks.txt"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create track tracker")
	}

	client.albumTracker, err = NewTracker(filepath.Join(baseDir, "albums.txt"))
	if err != nil {
		return nil, errors.Wrap(err, "unable to create album tracker")
	}

	return client, nil
}

func (client *Client) getBundle() (string, error) {
	if client.bundle != "" {
		return client.bundle, nil
	}

	bundleURL, err := client.getBundleURL()
	if err != nil {
		return "", errors.Wrap(err, "get bundle url")
	}

	// do some basic verification that the url is valid
	if !(strings.HasPrefix(bundleURL, "https://play.qobuz.com/resources/") &&
		strings.HasSuffix(bundleURL, "/bundle.js")) {
		return "", errors.New("invalid bundle url")
	}

	res, err := http.Get(bundleURL) //nolint:noctx,gosec // We do as much validation as we can above
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = errors.Wrap(err, "failed to close m3u file")
		}
	}()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return "", errors.Wrap(err, "read body")
	}

	client.bundle = string(buf)

	return client.bundle, nil
}

func (*Client) getBundleURL() (string, error) {
	res, err := http.Get(baseApp + "/login") //nolint:noctx
	if err != nil {
		return "", errors.Wrap(err, "do request")
	}

	defer func() {
		if closeErr := res.Body.Close(); closeErr != nil {
			err = errors.Wrap(err, "failed to close m3u file")
		}
	}()

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

func (client *Client) getAppID() (string, error) {
	bundle, err := client.getBundle()
	if err != nil {
		return "", errors.Wrap(err, "get bundle")
	}

	matches := appIDRegexp.FindAllStringSubmatch(bundle, -1)
	if len(matches) == 0 {
		return "", errors.New("no matches found")
	}

	return matches[0][1], nil
}

func (client *Client) testSecret(secret string) bool {
	secrets := client.Secrets
	client.Secrets = []string{secret}

	defer func() {
		client.Secrets = secrets
	}()

	_, err := client.TrackGetFileURL("5966783", QualityMP3)
	if err != nil && !errors.Is(err, common.ErrUnavailable) {
		if !errors.Is(err, common.ErrBadRequest) {
			log.Debug().Err(err).Msg("unexpected error when testing secrets")
		}

		return false
	}

	return true
}

//nolint:cyclop,funlen // todo: fix in future
func (client *Client) getSecrets() ([]string, error) {
	bundle, err := client.getBundle()
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

		if client.testSecret(string(base64Secret)) {
			secrets = append(secrets, string(base64Secret))
		}
	}

	return secrets, nil
}

func (client *Client) Login(email, password string) error {
	login, err := (Querier[userLogin.Response]{client}).Req("user/login", &url.Values{
		"email":    {email},
		"password": {password},
		"app_id":   {client.AppID},
	})
	if err != nil {
		return errors.Wrap(err, "user login")
	}

	if login.UserAuthToken == "" {
		return errors.New("no user auth token found")
	}

	client.Header.Set(userAuthToken, login.UserAuthToken)

	return nil
}

func (client *Client) Close() error {
	if err := client.trackTracker.Close(); err != nil {
		return errors.Wrap(err, "unable to close track tracker")
	}

	if err := client.albumTracker.Close(); err != nil {
		return errors.Wrap(err, "unable to close album tracker")
	}

	return nil
}

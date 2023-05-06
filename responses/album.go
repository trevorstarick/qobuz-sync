package responses

import (
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/trevorstarick/qobuz-sync/common"
	"github.com/trevorstarick/qobuz-sync/helpers"
)

//nolint:tagliatelle
type Album struct {
	Image                          Image     `json:"image"`
	MaximumBitDepth                int       `json:"maximum_bit_depth"`
	MediaCount                     int       `json:"media_count"`
	Artist                         *Artist   `json:"artist"`
	Artists                        []*Artist `json:"artists"`
	Upc                            string    `json:"upc"`
	ReleasedAt                     int       `json:"released_at"`
	Label                          *Artist   `json:"label"`
	Title                          string    `json:"title"`
	QobuzID                        int       `json:"qobuz_id"`
	URL                            string    `json:"url"`
	Version                        string    `json:"version"`
	Duration                       int       `json:"duration"`
	ParentalWarning                bool      `json:"parental_warning"`
	TracksCount                    int       `json:"tracks_count"`
	Popularity                     int       `json:"popularity"`
	Genre                          Genre     `json:"genre"`
	MaximumChannelCount            int       `json:"maximum_channel_count"`
	ID                             string    `json:"id"`
	MaximumSamplingRate            float64   `json:"maximum_sampling_rate"`
	Articles                       []any     `json:"articles"`
	Previewable                    bool      `json:"previewable"`
	Sampleable                     bool      `json:"sampleable"`
	Displayable                    bool      `json:"displayable"`
	Streamable                     bool      `json:"streamable"`
	StreamableAt                   int       `json:"streamable_at"`
	Downloadable                   bool      `json:"downloadable"`
	PurchasableAt                  any       `json:"purchasable_at"`
	Purchasable                    bool      `json:"purchasable"`
	ReleaseDateOriginal            string    `json:"release_date_original"`
	ReleaseDateDownload            string    `json:"release_date_download"`
	ReleaseDateStream              string    `json:"release_date_stream"`
	Hires                          bool      `json:"hires"`
	HiresStreamable                bool      `json:"hires_streamable"`
	FavoritedAt                    int       `json:"favorited_at"`
	Awards                         []any     `json:"awards"`
	Description                    string    `json:"description"`
	DescriptionLanguage            string    `json:"description_language"`
	Goodies                        []any     `json:"goodies"`
	Area                           any       `json:"area"`
	Catchline                      string    `json:"catchline"`
	Composer                       *Artist   `json:"composer"`
	CreatedAt                      int       `json:"created_at"`
	GenresList                     []string  `json:"genres_list"`
	Period                         any       `json:"period"`
	Copyright                      string    `json:"copyright"`
	IsOfficial                     bool      `json:"is_official"`
	MaximumTechnicalSpecifications string    `json:"maximum_technical_specifications"`
	ProductSalesFactorsMonthly     float64   `json:"product_sales_factors_monthly"`
	ProductSalesFactorsWeekly      float64   `json:"product_sales_factors_weekly"`
	ProductSalesFactorsYearly      float64   `json:"product_sales_factors_yearly"`
	ProductType                    string    `json:"product_type"`
	ProductURL                     string    `json:"product_url"`
	RecordingInformation           string    `json:"recording_information"`
	RelativeURL                    string    `json:"relative_url"`
	ReleaseTags                    []any     `json:"release_tags"`
	ReleaseType                    string    `json:"release_type"`
	Slug                           string    `json:"slug"`
	Subtitle                       string    `json:"subtitle"`
}

func (album *Album) Path() string {
	artist := helpers.SanitizeStringToPath(album.Artist.Name)
	albumName := helpers.SanitizeStringToPath(album.Title)

	return filepath.Join(artist, albumName)
}

func (album *Album) DownloadAlbumArt(dir string) error {
	err := os.MkdirAll(dir, common.DirPerm)
	if err != nil {
		return errors.Wrap(err, "failed to create dir")
	}

	_, err = os.Stat(filepath.Join(dir, "album.jpg"))
	if err == nil {
		return common.ErrAlreadyExists
	}

	url := strings.ReplaceAll(album.Image.Large, "_600", "_org")
	log.Debug().Str("url", url).Str("dir", dir).Msg("downloading album art")

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to get album art")
	}

	defer res.Body.Close()

	ff, err := os.OpenFile(filepath.Join(dir, "album.jpg"), os.O_CREATE|os.O_WRONLY, common.FilePerm)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
	}

	defer ff.Close()

	_, err = io.Copy(ff, res.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy response body")
	}

	return nil
}

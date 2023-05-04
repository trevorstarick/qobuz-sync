package TrackGet

type Response struct {
	MaximumBitDepth     int       `json:"maximum_bit_depth"`
	Copyright           string    `json:"copyright"`
	Performers          string    `json:"performers"`
	AudioInfo           AudioInfo `json:"audio_info"`
	Performer           Performer `json:"performer"`
	Album               Album     `json:"album"`
	Work                any       `json:"work"`
	Composer            Composer  `json:"composer,omitempty"`
	Isrc                string    `json:"isrc"`
	Title               string    `json:"title"`
	Version             any       `json:"version"`
	Duration            int       `json:"duration"`
	ParentalWarning     bool      `json:"parental_warning"`
	TrackNumber         int       `json:"track_number"`
	MaximumChannelCount int       `json:"maximum_channel_count"`
	ID                  int       `json:"id"`
	MediaNumber         int       `json:"media_number"`
	MaximumSamplingRate float64   `json:"maximum_sampling_rate"`
	Articles            []any     `json:"articles"`
	ReleaseDateOriginal any       `json:"release_date_original"`
	ReleaseDateDownload any       `json:"release_date_download"`
	ReleaseDateStream   any       `json:"release_date_stream"`
	Purchasable         bool      `json:"purchasable"`
	Streamable          bool      `json:"streamable"`
	Previewable         bool      `json:"previewable"`
	Sampleable          bool      `json:"sampleable"`
	Downloadable        bool      `json:"downloadable"`
	Displayable         bool      `json:"displayable"`
	PurchasableAt       any       `json:"purchasable_at"`
	StreamableAt        int       `json:"streamable_at"`
	Hires               bool      `json:"hires"`
	HiresStreamable     bool      `json:"hires_streamable"`
}
type AudioInfo struct {
	ReplaygainTrackGain float64 `json:"replaygain_track_gain"`
	ReplaygainTrackPeak float64 `json:"replaygain_track_peak"`
}
type Performer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Composer struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	AlbumsCount int    `json:"albums_count"`
	Picture     any    `json:"picture"`
	Image       any    `json:"image"`
}
type Image struct {
	Small     string `json:"small"`
	Thumbnail string `json:"thumbnail"`
	Large     string `json:"large"`
	Back      any    `json:"back"`
}
type Artist struct {
	Image       any    `json:"image"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	AlbumsCount int    `json:"albums_count"`
	Slug        string `json:"slug"`
	Picture     any    `json:"picture"`
}
type Artists struct {
	ID    int      `json:"id"`
	Name  string   `json:"name"`
	Roles []string `json:"roles"`
}
type Label struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	AlbumsCount int    `json:"albums_count"`
	SupplierID  int    `json:"supplier_id"`
	Slug        string `json:"slug"`
}
type Genre struct {
	Path  []int  `json:"path"`
	Color string `json:"color"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
}
type Album struct {
	MaximumBitDepth                int       `json:"maximum_bit_depth"`
	Image                          Image     `json:"image"`
	MediaCount                     int       `json:"media_count"`
	Artist                         Artist    `json:"artist"`
	Artists                        []Artists `json:"artists"`
	Upc                            string    `json:"upc"`
	ReleasedAt                     int       `json:"released_at"`
	Label                          Label     `json:"label"`
	Title                          string    `json:"title"`
	QobuzID                        int       `json:"qobuz_id"`
	Version                        any       `json:"version"`
	URL                            string    `json:"url"`
	Duration                       int       `json:"duration"`
	ParentalWarning                bool      `json:"parental_warning"`
	Popularity                     int       `json:"popularity"`
	TracksCount                    int       `json:"tracks_count"`
	Genre                          Genre     `json:"genre"`
	MaximumChannelCount            int       `json:"maximum_channel_count"`
	ID                             string    `json:"id"`
	MaximumSamplingRate            float64   `json:"maximum_sampling_rate"`
	Articles                       []any     `json:"articles"`
	ReleaseDateOriginal            string    `json:"release_date_original"`
	ReleaseDateDownload            string    `json:"release_date_download"`
	ReleaseDateStream              string    `json:"release_date_stream"`
	Purchasable                    bool      `json:"purchasable"`
	Streamable                     bool      `json:"streamable"`
	Previewable                    bool      `json:"previewable"`
	Sampleable                     bool      `json:"sampleable"`
	Downloadable                   bool      `json:"downloadable"`
	Displayable                    bool      `json:"displayable"`
	PurchasableAt                  int       `json:"purchasable_at"`
	StreamableAt                   int       `json:"streamable_at"`
	Hires                          bool      `json:"hires"`
	HiresStreamable                bool      `json:"hires_streamable"`
	Awards                         []any     `json:"awards"`
	Goodies                        []any     `json:"goodies"`
	Area                           any       `json:"area"`
	Catchline                      string    `json:"catchline"`
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
	Description                    string    `json:"description"`
}

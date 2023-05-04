package TrackSearch

type TrackSearch struct {
	Query  string `json:"query"`
	Tracks Tracks `json:"tracks"`
}
type Analytics struct {
	SearchExternalID string `json:"search_external_id"`
}
type AudioInfo struct {
	ReplaygainTrackPeak float64 `json:"replaygain_track_peak"`
	ReplaygainTrackGain float64 `json:"replaygain_track_gain"`
}
type Performer struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
type Image struct {
	Small     string `json:"small"`
	Thumbnail string `json:"thumbnail"`
	Large     string `json:"large"`
}
type Artist struct {
	Image       any    `json:"image"`
	Name        string `json:"name"`
	ID          int    `json:"id"`
	AlbumsCount int    `json:"albums_count"`
	Slug        string `json:"slug"`
	Picture     any    `json:"picture"`
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
	Image               Image   `json:"image"`
	MaximumBitDepth     int     `json:"maximum_bit_depth"`
	MediaCount          int     `json:"media_count"`
	Artist              Artist  `json:"artist"`
	Upc                 string  `json:"upc"`
	ReleasedAt          int     `json:"released_at"`
	Label               Label   `json:"label"`
	Title               string  `json:"title"`
	QobuzID             int     `json:"qobuz_id"`
	Version             any     `json:"version"`
	Duration            int     `json:"duration"`
	ParentalWarning     bool    `json:"parental_warning"`
	TracksCount         int     `json:"tracks_count"`
	Popularity          int     `json:"popularity"`
	Genre               Genre   `json:"genre"`
	MaximumChannelCount int     `json:"maximum_channel_count"`
	ID                  string  `json:"id"`
	MaximumSamplingRate float64 `json:"maximum_sampling_rate"`
	Previewable         bool    `json:"previewable"`
	Sampleable          bool    `json:"sampleable"`
	Displayable         bool    `json:"displayable"`
	Streamable          bool    `json:"streamable"`
	StreamableAt        int     `json:"streamable_at"`
	Downloadable        bool    `json:"downloadable"`
	PurchasableAt       any     `json:"purchasable_at"`
	Purchasable         bool    `json:"purchasable"`
	ReleaseDateOriginal string  `json:"release_date_original"`
	ReleaseDateDownload string  `json:"release_date_download"`
	ReleaseDateStream   string  `json:"release_date_stream"`
	Hires               bool    `json:"hires"`
	HiresStreamable     bool    `json:"hires_streamable"`
}
type Composer struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}
type Items struct {
	MaximumBitDepth     int       `json:"maximum_bit_depth"`
	Copyright           string    `json:"copyright"`
	Performers          string    `json:"performers"`
	AudioInfo           AudioInfo `json:"audio_info"`
	Performer           Performer `json:"performer"`
	Album               Album     `json:"album"`
	Work                any       `json:"work"`
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
	Composer            Composer  `json:"composer,omitempty"`
}
type Tracks struct {
	Limit     int       `json:"limit"`
	Offset    int       `json:"offset"`
	Analytics Analytics `json:"analytics"`
	Total     int       `json:"total"`
	Items     []Items   `json:"items"`
}

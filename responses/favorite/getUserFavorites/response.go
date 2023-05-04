package FavoritesGetUserFavorites

type Response struct {
	Albums  AlbumsRes  `json:"albums"`
	Tracks  TracksRes  `json:"tracks"`
	Artists ArtistsRes `json:"artists"`
	User    User       `json:"user"`
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
	Small      string `json:"small"`
	Thumbnail  string `json:"thumbnail"`
	Large      string `json:"large"`
	Extralarge string `json:"extralarge"`
	Mega       string `json:"mega"`
	Back       any    `json:"back"`
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
	Image               Image   `json:"image"`
	MaximumBitDepth     int     `json:"maximum_bit_depth"`
	MediaCount          int     `json:"media_count"`
	Artist              Artist  `json:"artist"`
	Upc                 string  `json:"upc"`
	ReleasedAt          int     `json:"released_at"`
	Label               Label   `json:"label"`
	Title               string  `json:"title"`
	QobuzID             int     `json:"qobuz_id"`
	Version             string  `json:"version"`
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
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	AlbumsCount int    `json:"albums_count"`
	Picture     any    `json:"picture"`
	Image       any    `json:"image"`
}
type TracksItems struct {
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
	Version             string    `json:"version"`
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
	PurchasableAt       int       `json:"purchasable_at"`
	StreamableAt        int       `json:"streamable_at"`
	Hires               bool      `json:"hires"`
	HiresStreamable     bool      `json:"hires_streamable"`
	FavoritedAt         int       `json:"favorited_at"`
}
type TracksRes struct {
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Total  int           `json:"total"`
	Items  []TracksItems `json:"items"`
}
type AlbumsItems struct {
	MaximumBitDepth     int       `json:"maximum_bit_depth"`
	Image               Image     `json:"image"`
	MediaCount          int       `json:"media_count"`
	Artist              Artist    `json:"artist"`
	Artists             []Artists `json:"artists"`
	Upc                 string    `json:"upc"`
	ReleasedAt          int       `json:"released_at"`
	Label               Label     `json:"label"`
	Title               string    `json:"title"`
	QobuzID             int       `json:"qobuz_id"`
	Version             any       `json:"version"`
	URL                 string    `json:"url"`
	Duration            int       `json:"duration"`
	ParentalWarning     bool      `json:"parental_warning"`
	Popularity          int       `json:"popularity"`
	TracksCount         int       `json:"tracks_count"`
	Genre               Genre     `json:"genre"`
	MaximumChannelCount int       `json:"maximum_channel_count"`
	ID                  string    `json:"id"`
	MaximumSamplingRate float64   `json:"maximum_sampling_rate"`
	Articles            []any     `json:"articles"`
	ReleaseDateOriginal string    `json:"release_date_original"`
	ReleaseDateDownload string    `json:"release_date_download"`
	ReleaseDateStream   string    `json:"release_date_stream"`
	Purchasable         bool      `json:"purchasable"`
	Streamable          bool      `json:"streamable"`
	Previewable         bool      `json:"previewable"`
	Sampleable          bool      `json:"sampleable"`
	Downloadable        bool      `json:"downloadable"`
	Displayable         bool      `json:"displayable"`
	PurchasableAt       int       `json:"purchasable_at"`
	StreamableAt        int       `json:"streamable_at"`
	Hires               bool      `json:"hires"`
	HiresStreamable     bool      `json:"hires_streamable"`
	FavoritedAt         int       `json:"favorited_at"`
}
type AlbumsRes struct {
	Offset int           `json:"offset"`
	Limit  int           `json:"limit"`
	Total  int           `json:"total"`
	Items  []AlbumsItems `json:"items"`
}
type ArtistsItems struct {
	Name        string `json:"name"`
	ID          int    `json:"id"`
	AlbumsCount int    `json:"albums_count"`
	Slug        string `json:"slug"`
	Picture     any    `json:"picture"`
	Image       Image  `json:"image"`
	FavoritedAt int    `json:"favorited_at"`
}
type ArtistsRes struct {
	Offset int            `json:"offset"`
	Limit  int            `json:"limit"`
	Total  int            `json:"total"`
	Items  []ArtistsItems `json:"items"`
}
type User struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
}

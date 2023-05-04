package GetFileUrl

type Response struct {
	TrackID      int            `json:"track_id"`
	Duration     int            `json:"duration"`
	URL          string         `json:"url"`
	FormatID     int            `json:"format_id"`
	MimeType     string         `json:"mime_type"`
	Restrictions []Restrictions `json:"restrictions"`
	SamplingRate float64        `json:"sampling_rate"`
	BitDepth     int            `json:"bit_depth"`
}
type Restrictions struct {
	Code string `json:"code"`
}

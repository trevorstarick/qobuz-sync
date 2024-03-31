package tracksearch

import "github.com/trevorstarick/qobuz-sync/responses"

type TrackSearch struct {
	Query  string `json:"query"`
	Tracks struct {
		Limit     int `json:"limit"`
		Offset    int `json:"offset"`
		Analytics struct {
			SearchExternalID string `json:"search_external_id"`
		} `json:"analytics"`
		Total int               `json:"total"`
		Items []responses.Track `json:"items"`
	} `json:"tracks"`
}

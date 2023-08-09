package catalogsearch

import "github.com/trevorstarick/qobuz-sync/responses"
import "github.com/trevorstarick/qobuz-sync/responses/playlist/get"

// CatalogSearch is the response to a catalog/search request
type CatalogSearch struct {
	Query string `json:"query"`

	Albums struct {
		Limit     int `json:"limit"`
		Offset    int `json:"offset"`
		Analytics struct {
			SearchExternalID string `json:"search_external_id"`
		} `json:"analytics"`
		Total int               `json:"total"`
		Items []responses.Album `json:"items"`
	} `json:"albums"`

	Artists struct {
		Limit     int `json:"limit"`
		Offset    int `json:"offset"`
		Analytics struct {
			SearchExternalID string `json:"search_external_id"`
		} `json:"analytics"`
		Total int                `json:"total"`
		Items []responses.Artist `json:"items"`
	} `json:"artists"`

	Playlists struct {
		Limit     int `json:"limit"`
		Offset    int `json:"offset"`
		Analytics struct {
			SearchExternalID string `json:"search_external_id"`
		} `json:"analytics"`
		Total int                    `json:"total"`
		Items []playlistget.Response `json:"items"`
	} `json:"playlists"`

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

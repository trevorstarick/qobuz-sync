//nolint:misspell
package userlogin

type Response struct {
	User          User   `json:"user"`
	UserAuthToken string `json:"user_auth_token"`
}
type Subscription struct {
	Offer            string `json:"offer"`
	Periodicity      string `json:"periodicity"`
	StartDate        string `json:"start_date"`
	EndDate          string `json:"end_date"`
	IsCanceled       bool   `json:"is_canceled"`
	HouseholdSizeMax int    `json:"household_size_max"`
}
type ColorScheme struct {
	Logo string `json:"logo"`
}

//nolint:revive
type Parameters struct {
	LossyStreaming          bool        `json:"lossy_streaming"`
	LosslessStreaming       bool        `json:"lossless_streaming"`
	HiresStreaming          bool        `json:"hires_streaming"`
	HiresPurchasesStreaming bool        `json:"hires_purchases_streaming"`
	MobileStreaming         bool        `json:"mobile_streaming"`
	OfflineStreaming        bool        `json:"offline_streaming"`
	HfpPurchase             bool        `json:"hfp_purchase"`
	IncludedFormatGroupIds  []int       `json:"included_format_group_ids"`
	ColorScheme             ColorScheme `json:"color_scheme"`
	Label                   string      `json:"label"`
	ShortLabel              string      `json:"short_label"`
	Source                  string      `json:"source"`
}
type Credential struct {
	ID          int        `json:"id"`
	Label       string     `json:"label"`
	Description string     `json:"description"`
	Parameters  Parameters `json:"parameters"`
}
type LastUpdate struct {
	Favorite       int `json:"favorite"`
	FavoriteAlbum  int `json:"favorite_album"`
	FavoriteArtist int `json:"favorite_artist"`
	FavoriteTrack  int `json:"favorite_track"`
	Playlist       int `json:"playlist"`
	Purchase       int `json:"purchase"`
}
type StoreFeatures struct {
	Download                 bool `json:"download"`
	Streaming                bool `json:"streaming"`
	Editorial                bool `json:"editorial"`
	Club                     bool `json:"club"`
	Wallet                   bool `json:"wallet"`
	Weeklyq                  bool `json:"weeklyq"`
	Autoplay                 bool `json:"autoplay"`
	InappPurchaseSubscripton bool `json:"inapp_purchase_subscripton"`
	OptIn                    bool `json:"opt_in"`
	MusicImport              bool `json:"music_import"`
}
type PlayerSettings struct {
	SonosAudioFormat int `json:"sonos_audio_format"`
}
type Settings struct {
	Scrobbling string `json:"scrobbling"`
}
type Lastfm struct {
	Name     string   `json:"name"`
	Key      string   `json:"key"`
	Settings Settings `json:"settings"`
}
type Externals struct {
	Lastfm Lastfm `json:"lastfm"`
}
type User struct {
	ID             int            `json:"id"`
	PublicID       string         `json:"publicId"`
	Email          string         `json:"email"`
	Login          string         `json:"login"`
	Firstname      any            `json:"firstname"`
	Lastname       any            `json:"lastname"`
	DisplayName    string         `json:"display_name"`
	CountryCode    string         `json:"country_code"`
	LanguageCode   string         `json:"language_code"`
	Zone           string         `json:"zone"`
	Store          string         `json:"store"`
	Country        string         `json:"country"`
	Avatar         string         `json:"avatar"`
	Genre          string         `json:"genre"`
	Age            int            `json:"age"`
	CreationDate   string         `json:"creation_date"`
	Subscription   Subscription   `json:"subscription"`
	Credential     Credential     `json:"credential"`
	LastUpdate     LastUpdate     `json:"last_update"`
	StoreFeatures  StoreFeatures  `json:"store_features"`
	PlayerSettings PlayerSettings `json:"player_settings"`
	Externals      Externals      `json:"externals"`
}

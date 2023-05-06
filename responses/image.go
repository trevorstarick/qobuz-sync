package responses

type Image struct {
	Small      string `json:"small"`
	Thumbnail  string `json:"thumbnail"`
	Large      string `json:"large"`
	Extralarge string `json:"extralarge"`
	Mega       string `json:"mega"`
	Back       any    `json:"back"`
}

package responses

type Genre struct {
	Path  []int  `json:"path"`
	Color string `json:"color"`
	Name  string `json:"name"`
	ID    int    `json:"id"`
	Slug  string `json:"slug"`
}

package api

type Landa struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	URL  string `json:"-"` // Maybe call it service or whatever in the future
}

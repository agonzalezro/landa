package model

type Landa struct {
	ID         string `json:"id"`
	Code       string `json:"code"`
	EntryPoint string `json:"entryPoint"`
	URL        string `json:"-"`
}

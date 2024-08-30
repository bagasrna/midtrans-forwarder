package domain

type Reseller struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Code  string `json:"code"`
	URL   string `json:"url"`
	Token string `json:"token"`
}
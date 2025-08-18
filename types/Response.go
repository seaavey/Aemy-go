package types

import "net/http"


type ResponseAPIs struct {
	Status int
	Body []byte
	Headers http.Header
}

type Image struct {
	URL string `json:"url"`
	Width int `json:"width"`
	Height int `json:"height"`
}

type TiktokResponse struct {
	Status int `json:"status"`
	Data struct {
		Title string `json:"title"`
		Video *struct {
			NoWatermark string `json:"noWatermark"`
		} `json:"video"`
		Music struct {
			PlayURL string `json:"play_url"`
		} `json:"music"`
		Images []Image `json:"images"`
	} `json:"data"`
}
type InstagramResponse struct {
	Creator string   `json:"creator"`
	Status  int      `json:"status"`
	Data    []string `json:"data"`
}
package types

type TiktokResponse struct {
	Status int `json:"status"`
	Data   struct {
		Title string `json:"title"`
		Video struct {
			NoWatermark string `json:"noWatermark"`
		} `json:"video"`
		Music struct {
			PlayURL string `json:"play_url"`
		} `json:"music"`
	} `json:"data"`
}

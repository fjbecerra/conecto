package rest


type ResponseConfig struct {
	BaseUrl string `json:"base_url"`
	Data struct {
		Path    string `json:"path"`
		IsArray bool   `json:"is_array"`
	}`json:"data"`

	Pagination struct {
		Type string `json:"type"`

		Request struct {
			Param string `json:"param"`
		} `json:"request"`

		Response struct {
			Next struct {
				Path string `json:"path"`
			} `json:"next"`
		} `json:"response"`
	} `json:"pagination"`
}

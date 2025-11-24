package model

type MetaInfo struct {
	CurrentPage int    `json:"current_page"`
	Limit       int    `json:"limit"`
	Total       int    `json:"total"`
	Pages       int    `json:"pages"`
	SortBy      string `json:"sort_by"`
	Order       string `json:"order"`
	Search      string `json:"search"`
}

type UserResponse struct {
	Data     []Users  `json:"data"`
	MetaInfo MetaInfo `json:"meta_info"`
}
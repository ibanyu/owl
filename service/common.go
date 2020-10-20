package service

type Pagination struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Key    string `json:"key"`

	Operator string `json:"operator"`
}


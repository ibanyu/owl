package controller

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListData struct {
	Total  int         `json:"total"`
	Items  interface{} `json:"items"`
	More   bool        `json:"more"`
	Offset int         `json:"offset"`
}

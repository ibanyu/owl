package controller

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListData struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

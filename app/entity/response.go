package entity

type HttpResponse struct {
	Timestamp  int64
	StatusCode int
	Body       string
}

type HttpResult struct {
	Response HttpResponse
	Error    error
}

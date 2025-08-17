package entity

type HttpResponse struct {
	Timestamp  int64
	StatusCode int
	Body       string
	Header     map[string][]string
}

type HttpResult struct {
	Response HttpResponse
	Error    error
}

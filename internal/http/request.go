package http


type HttpRequestMethod string

const (
	GET    HttpRequestMethod = "GET"
	POST   HttpRequestMethod = "POST"
	PUT    HttpRequestMethod = "PUT"
	DELETE HttpRequestMethod = "DELETE"
	PATCH  HttpRequestMethod = "PATCH"
)

type HttpRequest struct {
	Url         string
	Headers     map[string]string
	Method      HttpRequestMethod
	Body        string
}

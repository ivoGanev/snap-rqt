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

func GetTcellColorForRequest(method HttpRequestMethod) string {
	switch method {
	case GET:
		return "[#942f94]" // Purple (similar to CSS .get)
	case POST:
		return "[green]" // Green (CSS .post)
	case PUT:
		return "[orange]" // Orange (CSS .put)
	case PATCH:
		return "[#a7a157]" // Yellowish (CSS .patch)
	case DELETE:
		return "[#d82929]" // Red (CSS .delete)
	default:
		return "[white]" // Default color for unknown methods
	}
}
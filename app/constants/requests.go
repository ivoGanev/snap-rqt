package constants

type RequestMethod string

const (
	GET    RequestMethod = "GET"
	POST   RequestMethod = "POST"
	PUT    RequestMethod = "PUT"
	DELETE RequestMethod = "DELETE"
	PATCH  RequestMethod = "PATCH"
)

var RequestMethods = []RequestMethod{GET, POST, PUT, DELETE, PATCH}

func RequestMethodStrings() []string {
	result := make([]string, len(RequestMethods))
	for i, m := range  RequestMethods {
		result[i] = string(m)
	}
	return result
}

func RequestMethodIndex(method RequestMethod) int {
	for i, m := range RequestMethods {
		if m == method {
			return i
		}
	}
	return -1 // not found
}

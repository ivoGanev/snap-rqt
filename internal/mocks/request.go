package mocks

import (
	"math/rand"
	"snap-rq/internal/data"
	"snap-rq/internal/http"

	"github.com/google/uuid"
)

var sampleRequests = map[string]struct {
	description string
	url         string
}{
	"Users":    {"user details", "https://api.example.com/users"},
	"Orders":   {"orders", "https://api.example.com/orders"},
	"Accounts": {"account", "https://api.example.com/account"},
	"Products": {"product", "https://api.example.com/products"},
	"Profiles": {"profile settings", "https://api.example.com/profile"},
}

var httpMethods = []http.RequestMethod{http.GET, http.POST, http.PUT, http.DELETE, http.PATCH}

func generateRandomHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	if rand.Intn(2) == 1 {
		headers["Authorization"] = "Bearer " + uuid.New().String()
	}
	return headers
}

func generateRandomBody() string {
	bodies := []string{
		`{"key": "value"}`,
		`{"username": "user123", "password": "securepass"}`,
		`{"product": "Laptop", "price": 999.99}`,
		`{"order_id": "12345", "status": "shipped"}`,
		`{"email": "example@email.com", "verified": true}`,
	}
	return bodies[rand.Intn(len(bodies))]
}

func GenerateMockRequests(count int) *[]data.Node[http.Request] {
	var nodes []data.Node[http.Request]

	keys := make([]string, 0, len(sampleRequests))
	for key := range sampleRequests {
		keys = append(keys, key)
	}

	for i := 0; i < count; i++ {
		key := keys[rand.Intn(len(keys))]
		entry := sampleRequests[key]
		method := httpMethods[rand.Intn(len(httpMethods))]
		headers := generateRandomHeaders()

		var body string
		if method == http.POST || method == http.PUT || method == http.PATCH {
			body = generateRandomBody()
		}

		request := &http.Request{
			Url:     entry.url,
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		node := data.NewNode(key, entry.description, request)
		nodes = append(nodes, node)
	}

	return &nodes
}

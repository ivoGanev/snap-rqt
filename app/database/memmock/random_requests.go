package memmock

import (
	"math/rand"
	"snap-rq/app"
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

func GenerateMockRequests(count int, collectionId string) []app.Request {
	var nodes []app.Request

	keys := make([]string, 0, len(sampleRequests))
	for key := range sampleRequests {
		keys = append(keys, key)
	}

	for range count {
		name := keys[rand.Intn(len(keys))]
		entry := sampleRequests[name]
		method := app.RequestMethods[rand.Intn(len(app.RequestMethods))]
		headers := generateRandomHeaders()

		var body string
		if method == "POST" || method == "PUT" || method == "PATCH" {
			body = generateRandomBody()
		}

		request := app.NewRequest(collectionId, name, entry.description, string(method), entry.url, headers, body)

		nodes = append(nodes, request)
	}

	return nodes
}

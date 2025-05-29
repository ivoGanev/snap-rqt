package random

import (
	"math/rand"
	"snap-rq/app/constants"
	"snap-rq/app/entity"

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

func Requests(count int, collectionId string) []entity.Request {
	var nodes []entity.Request

	keys := make([]string, 0, len(sampleRequests))
	for key := range sampleRequests {
		keys = append(keys, key)
	}

	for index := range count {
		name := keys[rand.Intn(len(keys))]
		entry := sampleRequests[name]
		method := constants.RequestMethods[rand.Intn(len(constants.RequestMethods))]
		headers := generateRandomHeaders()

		var body string
		if method == "POST" || method == "PUT" || method == "PATCH" {
			body = generateRandomBody()
		}

		request := entity.NewRequest(collectionId, name, entry.description, string(method), entry.url, headers, body, index)

		nodes = append(nodes, request)
	}

	return nodes
}

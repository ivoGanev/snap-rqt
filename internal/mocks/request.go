package mocks

import (
	"math/rand"
	"snap-rq/internal/data"
	"snap-rq/internal/http"

	"github.com/google/uuid"
)

// Mock API Endpoints
var apiEndpoints = []string{
	"https://api.example.com/users",
	"https://api.example.com/orders",
	"https://api.example.com/account",
	"https://api.example.com/products",
	"https://api.example.com/profile",
}

var httpMethods = []http.RequestMethod{http.GET, http.POST, http.PUT, http.DELETE, http.PATCH}

var sampleNames = []string{"Users That Don't Make Any Sense", "Orders With Big Gains",  "Accounts", "Products", "Profiles"}
var sampleDescriptions = []string{
	"Fetch user details",
	"Update orders",
	"Delete an account",
	"Create a new product",
	"Modify profile settings",
}

func generateRandomHeaders() map[string]string {
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// Randomly add an Authorization header
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

func GenerateMockRequests(count int) []data.Node[http.Request] {
	var nodes []data.Node[http.Request]

	for i := 0; i < count; i++ {
		randomNumber := rand.Intn(len(sampleNames))
		name := sampleNames[randomNumber]
		description := sampleDescriptions[randomNumber]
		url := apiEndpoints[randomNumber]
		method := httpMethods[rand.Intn(len(httpMethods))]
		headers := generateRandomHeaders()

		var body string
		if method == http.POST || method == http.PUT || method == http.PATCH {
			body = generateRandomBody()
		}

		request := &http.Request{
			Url:     url,
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		node := data.NewNode(name, description, request)
		nodes = append(nodes, node)
	}

	return nodes
}

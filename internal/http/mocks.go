package http

import (
	"math/rand"
	"snap-rq/internal/data"
	"time"

	"github.com/google/uuid"
)

// Possible API Endpoints
var apiEndpoints = []string{
	"https://api.example.com/users",
	"https://api.example.com/posts",
	"https://api.example.com/comments",
	"https://api.example.com/products",
	"https://api.example.com/orders",
}

// Possible HTTP Methods
var httpMethods = []HttpRequestMethod{GET, POST, PUT, DELETE, PATCH}

// Sample Names and Descriptions
var sampleNames = []string{"Users That Don't Make Any Sense", "Orders With Big Gains", "Products", "Accounts", "Profiles"}
var sampleDescriptions = []string{
	"Fetch user details",
	"Update order information",
	"Create a new product",
	"Delete an account",
	"Modify profile settings",
}

// Generate Random Headers
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

// Generate Random JSON Body (for POST, PUT, PATCH)
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

// GenerateMockRequestsAsNodes - Creates multiple random API requests
func GenerateMockRequestsAsNodes(count int) []data.Node[HttpRequest] {
	var nodes []data.Node[HttpRequest]
	rand.Seed(time.Now().UnixNano()) // Seed the random generator

	for i := 0; i < count; i++ {
		name := sampleNames[rand.Intn(len(sampleNames))]
		description := sampleDescriptions[rand.Intn(len(sampleDescriptions))]
		url := apiEndpoints[rand.Intn(len(apiEndpoints))]
		method := httpMethods[rand.Intn(len(httpMethods))]
		headers := generateRandomHeaders()

		var body string
		if method == POST || method == PUT || method == PATCH {
			body = generateRandomBody()
		}

		request := &HttpRequest{
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
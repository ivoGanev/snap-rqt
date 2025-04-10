package data

import (
	"math/rand"
	"snap-rq/internal/http"

	"github.com/google/uuid"
)

type MockStore struct {
	StoredCollections []Collection
	StoredRequests    []Request
}

func (m *MockStore) StoreCollection(c Collection) error {
	m.StoredCollections = append(m.StoredCollections, c)
	return nil
}

func (m *MockStore) StoreRequest(r Request) error {
	m.StoredRequests = append(m.StoredRequests, r)
	return nil
}

func (m *MockStore) LoadAllCollections() ([]Collection, error) {
	return m.StoredCollections, nil
}

func (m *MockStore) LoadAllRequests() ([]Request, error) {
	return m.StoredRequests, nil
}

func (m *MockStore) GetCollectionsForRequest(requestID int) ([]Collection, error) {
	// simple mock: return all collections for simplicity
	return m.StoredCollections, nil
}




var collectionNames = []string{
	"User APIs",
	"Order Services",
	"Product Endpoints",
	"Authentication Flows",
	"Billing & Payments",
	"Admin Tools",
	"Notification Services",
	"Analytics APIs",
	"Third-Party Integrations",
	"Debugging & Testing APIs",
}

func GenerateCollectionMocks(collectionCount int, requestGenerator func() *[]Node[http.Request]) *[]Collection {
	var collections []Collection

	for range collectionCount {
		collectionName := collectionNames[rand.Intn(len(collectionNames))]
		collectionDescription := "A collection of API requests for " + collectionName
		collection := Collection { Node: NewNode(collectionName, collectionDescription, requestGenerator()) }
		collections = append(collections, collection)
	}

	return &collections
}


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

func GenerateMockRequests(count int) *[]Node[http.Request] {
	var nodes []Node[http.Request]

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

		node := NewNode(key, entry.description, request)
		nodes = append(nodes, node)
	}

	return &nodes
}

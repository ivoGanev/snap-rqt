package http
import (
	"snap-rq/internal/data"
)

func GetMockRequestsAsNodes() []data.Node[HttpRequest] {
	return []data.Node[HttpRequest]{
		data.NewNode("Get Users", "Fetch all users from API", &HttpRequest{
			Url:    "https://api.example.com/users",
			Method: GET,
			Headers: map[string]string{}, // Empty headers
		}),
		data.NewNode("User Login", "Authenticate a user", &HttpRequest{
			Url:    "https://api.example.com/login",
			Method: POST,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"username": "test", "password": "1234"}`,
		}),
		data.NewNode("Update User", "Modify user data", &HttpRequest{
			Url:    "https://api.example.com/update",
			Method: PUT,
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
			Body: `{"key": "newValue"}`,
		}),
		data.NewNode("Delete Account", "Remove a user account", &HttpRequest{
			Url:    "https://api.example.com/delete",
			Method: DELETE,
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
		}),
		data.NewNode("Update Profile", "Modify user profile data", &HttpRequest{
			Url:    "https://api.example.com/profile",
			Method: PATCH,
			Headers: map[string]string{
				"X-Request-ID": "12345",
			},
			Body: `{"profile": "updated"}`,
		}),
	}
}
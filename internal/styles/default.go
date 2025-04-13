package styles

import (
	"fmt"
)

type Default struct {

}

func (s *Default) GetStyledRequestMethod(method string) string {
	var color string

	switch method {
	case "GET":
		color = "[#942f94]" // Purple
	case "POST":
		color = "[green]"
	case "PUT":
		color = "[#ffa500]"
	case "PATCH":
		color = "[#a7a157]" // Brownish
	case "DELETE":
		color = "[#d82929]" // Red
	default:
		color = "[white]"
	}

	return fmt.Sprintf("%s %s [white]", color, method)
}

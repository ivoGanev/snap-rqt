package styles

type StyleProvider interface {
	GetStyledRequestMethod(string) string
}
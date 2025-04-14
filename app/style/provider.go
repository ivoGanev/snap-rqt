package style

type StyleProvider interface {
	GetStyledRequestMethod(string) string
}

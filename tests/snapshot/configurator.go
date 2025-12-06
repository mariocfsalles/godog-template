package snapshot

type Config struct {
	File      string                                          // JSON file name, e.g. "store.json"
	CheckFunc func(fixturePath string, body []byte) error     // function that performs the assertion
}

// generic helper: creates a Config for T type
func NewConfig[T any](file string, normalize func(*T)) Config {
	return Config{
		File: file,
		CheckFunc: func(path string, body []byte) error {
			return AssertSnapshot(path, body, normalize)
		},
	}
}

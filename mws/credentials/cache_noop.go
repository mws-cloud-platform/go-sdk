package credentials

// NoopCache is a credentials cache implementation that does nothing.
type NoopCache struct {
}

// NewNoopCache creates a new [NoopCache] instance.
func NewNoopCache() *NoopCache {
	return new(NoopCache)
}

func (n *NoopCache) Load(string) (Credentials, error) {
	return Credentials{}, ErrEntryNotFound
}

func (n *NoopCache) Store(string, Credentials) error {
	return nil
}

func (n *NoopCache) Delete(string) error {
	return nil
}

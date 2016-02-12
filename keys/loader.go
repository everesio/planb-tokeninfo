package keys

type KeyLoader interface {
	LoadKey(id string) (interface{}, error)
}

func DefaultKeyLoader() KeyLoader {
	return newCachingOpenIdProviderLoader()
}
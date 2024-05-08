package storage

type URLMapper interface {
	Get(shortLink string) (string, bool)
	Set(shortLink, longLink string)
}

type SimpleURLMapper struct {
	URLs map[string]string
}

func newSimpleURLMapper() *SimpleURLMapper {
	return &SimpleURLMapper{URLs: make(map[string]string)}
}

func (m *SimpleURLMapper) Get(shortLink string) (string, bool) {
	longLink, ok := m.URLs[shortLink]
	return longLink, ok
}

func (m *SimpleURLMapper) Set(shortLink, longLink string) {
	m.URLs[shortLink] = longLink
}

var Mapper = newSimpleURLMapper()

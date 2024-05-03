package storage

type SimpleURLMapper struct {
	URLs map[string]string
}

func NewSimpleURLMapper() *SimpleURLMapper {
	return &SimpleURLMapper{URLs: make(map[string]string)}
}

func (m *SimpleURLMapper) Get(shortLink string) (string, bool) {
	longLink, ok := m.URLs[shortLink]
	return longLink, ok
}

func (m *SimpleURLMapper) Set(shortLink, longLink string) {
	m.URLs[shortLink] = longLink
}

func (m *SimpleURLMapper) Count() int {
	return len(m.URLs)
}

var Mapper = NewSimpleURLMapper()

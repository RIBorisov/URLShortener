package storage

type SimpleURLMapper struct {
	urls map[string]string
}

func NewSimpleURLMapper() *SimpleURLMapper {
	return &SimpleURLMapper{urls: make(map[string]string)}
}

func (m *SimpleURLMapper) Get(shortLink string) (string, bool) {
	longLink, ok := m.urls[shortLink]
	return longLink, ok
}

func (m *SimpleURLMapper) Set(shortLink, longLink string) {
	m.urls[shortLink] = longLink
}

func (m *SimpleURLMapper) Count() int {
	return len(m.urls)
}

var Mapper = NewSimpleURLMapper()

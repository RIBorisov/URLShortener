package storage

type Storage struct {
	URLs map[string]string
}

func (m *Storage) Get(shortLink string) (string, bool) {
	longLink, ok := m.URLs[shortLink]
	return longLink, ok
}

func (m *Storage) Save(shortLink, longLink string) {
	m.URLs[shortLink] = longLink
}

func LoadStorage() *Storage {
	return &Storage{URLs: make(map[string]string)}
}

package storage

type Storage interface {
	Get(shortLink string) (string, bool)
	Save(shortLink, longLink string)
}

type SimpleStorage struct {
	URLs map[string]string
}

func (m *SimpleStorage) Get(shortLink string) (string, bool) {
	longLink, ok := m.URLs[shortLink]
	return longLink, ok
}

func (m *SimpleStorage) Save(shortLink, longLink string) {
	m.URLs[shortLink] = longLink
}

func GetStorage() *SimpleStorage {
	return &SimpleStorage{URLs: make(map[string]string)}
}

package storage

type URLMapper interface {
	Get(shortLink string) (string, bool)
	Set(shortLink, longLink string)
	Count() int
	// Delete() Not implemented yet
}

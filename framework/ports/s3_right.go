package ports

type S3Storage interface {
	WriteData(key, data string) error
	FetchData(key string) (string, error)
}

package cache

type CacheInterface interface {
	Get(string) ([]byte, error)
	Set(string, []byte) error
	Del(string) error
}

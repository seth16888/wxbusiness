package captcha

type Store interface {
	Set(key string, value string) error
	Get(key string, clear bool) string
	Verify(key, answer string, clear bool) bool
}

type MemoryStore struct {
	// 存储验证码的 map
	store map[string]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{store: make(map[string]string)}
}

func (c *MemoryStore) Set(key string, value string) error {
	c.store[key] = value
	return nil
}

func (c *MemoryStore) Get(key string, clear bool) string {
	if clear {
		value := c.store[key]
		delete(c.store, key)
		return value
	}
	return c.store[key]
}

func (c *MemoryStore) Verify(key, answer string, clear bool) bool {
	return c.Get(key, clear) == answer
}

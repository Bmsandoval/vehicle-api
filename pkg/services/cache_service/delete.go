package cache_service

func (s ServiceImplementation) Delete(key string) {
	s.Cache.Delete(key)
}

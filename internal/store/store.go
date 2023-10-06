package store

// Factory .
type Factory interface {
	Users() UserStore
}

// CacheFactory .
type CacheFactory interface {
	Users() UserCache
	Auth() AuthCache
}

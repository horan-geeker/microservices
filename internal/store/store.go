package store

// DataFactory .
type DataFactory interface {
	Users() UserStore
}

// CacheFactory .
type CacheFactory interface {
	Users() UserCache
	Auth() AuthCache
}

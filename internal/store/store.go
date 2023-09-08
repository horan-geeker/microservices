package store

// Factory .
type Factory interface {
	Users() UserStore
}

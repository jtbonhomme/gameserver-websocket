package storage

// Storage defines the interface for game manager storage.
type Storage interface {
	Players // Players defines the interface for players storage.
	Games   // Games defines the interface for games storage.
}

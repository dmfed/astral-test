package storage

// Storage defines basic CRUD interface
type Storage interface {
	// Get returns all existing elements if
	// called without arguments. If arguments
	// are present only the first ID is used.
	Get(...ID) ([]Element, error)
	// Put pushes Element to permanent storage
	// and returns an ID of the element.
	Put(Element) (ID, error)
	// Upd replaces Element with ID with the provided
	// Element.
	Upd(ID, Element) error
	// Del deletes Element with provided ID
	Del(ID) error
	// Close is intended to gracefully shitdown
	// Storage: disconnect from a DB, flush file to disk etc.
	Close() error
}

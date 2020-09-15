package models

// A Location represents any location at the school that can be signed in or out
type Location struct {
	ID uint64
	Name string
}

// GetLocations returns all of the locations stored in the database
func (db *Database) GetLocations() []Location {
	return nil
}
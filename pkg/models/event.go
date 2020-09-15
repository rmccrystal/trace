package models

// An Event represents a student either entering or leaving a location
type Event struct {
	ID uint64
	Location Location
}

// GetEvents gets all of the events in the database
func (db *Database) GetEvents() []Event {
	return nil
}
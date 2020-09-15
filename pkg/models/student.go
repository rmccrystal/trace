package models

// A Student represents one member of the school who can sign in and out of a location
type Student struct {
	ID uint64
	Name string
	Email string
}

// GetStudents gets a list of all students stored in the database
func (db *Database) GetStudents() []Student {
	return nil
}

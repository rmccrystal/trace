package database

import "go.mongodb.org/mongo-driver/bson/primitive"

// A Student represents one member of the school who can sign in and out of a location
type Student struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`
}

// GetStudents gets a list of all students stored in the Database
func (db *Database) GetStudents() []Student {
	return nil
}

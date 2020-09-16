package database

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A Student represents one member of the school who can sign in and out of a location
type Student struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `json:"name"`
	Email      string             `json:"email"`

	// StudentIDs is the list of IDs that can be used to scan in and out of a location
	StudentIDs []string           `json:"student_ids"`
}

// CreateStudent creates a student and adds it to the database. The
// ID element of the newly created Student will be set if it is successful
func (db *Database) CreateStudent(student *Student) error {
	result, err := db.Collections.Students.InsertOne(context.TODO(), student)
	if err != nil {
		return err
	}

	student.ID = result.InsertedID.(primitive.ObjectID)

	return nil
}

// GetStudents returns a list of all students stored in the database.
func (db *Database) GetStudents() ([]Student, error) {
	cur, err := db.Collections.Students.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}

	var students []Student
	if err := cur.All(context.TODO(), &students); err != nil {
		return nil, err
	}

	return students, nil
}

// GetStudentByStudentID gets a student by the StudentIDs member. If the
// student is found, found will be true. If there is an error getting the student,
// it will be returned. Not that this error will not contain the not found error
func (db *Database) GetStudentByStudentID(studentID string) (student Student, found bool, err error) {
	result := db.Collections.Students.FindOne(context.TODO(), bson.M{"studentids": bson.M{"$elemMatch": bson.M{"$eq": studentID}}})

	err = result.Err()
	if err != nil {
		// If the student cannot be found
		if err == mongo.ErrNoDocuments {
			return Student{}, false, nil
		}
		return Student{}, false, err
	}
	
	if err := result.Decode(&student); err != nil {
		return Student{}, false, err
	}

	found = true

	return
}

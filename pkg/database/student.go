package database

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// A Student represents one member of the school who can sign in and out of a location
type Student struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name  string             `json:"name"`
	Email string             `json:"email"`

	// StudentHandles is the list of IDs that can be used to scan in and out of a location
	StudentHandles []string `json:"student_handles"`
}

// StudentRef is a reference to a student which, when serialized, will return
// the json of the referenced student.
//
// Be careful for circular references.
type StudentRef primitive.ObjectID

func (ref StudentRef) GetBSON() (interface{}, error){
	return primitive.ObjectID(ref), nil
}

func (ref StudentRef) MarshalJSON() ([]byte, error) {
	obj, found := DB.GetStudentByID(primitive.ObjectID(ref))
	if !found {
		return nil, fmt.Errorf("could not find student with id %s", ref)
	}

	return json.Marshal(obj)
}

// Same functionality is ObjectID.UnmarshalJSON except it returns an error if the referenced object doesn't exist
func (ref *StudentRef) UnmarshalJSON(b []byte) error {
	id := primitive.ObjectID(*ref)
	if err := id.UnmarshalJSON(b); err != nil {
		return err
	}
	_, found := DB.GetEventByID(id)
	if !found {
		return fmt.Errorf("object not found")
	}

	*ref = StudentRef(id)
	return nil
}

// Gets the referenced object and panics if it doesn't exist
func (ref StudentRef) Get() Student {
	obj, found := DB.GetStudentByID(primitive.ObjectID(ref))
	if !found {
		panic("could not find object")
	}
	return obj
}

// Ref creates a reference to the object
func (obj Student) Ref() StudentRef {
	return StudentRef(obj.ID)
}

// CreateStudent creates a student and adds it to the database. The
// ID element of the newly created Student will be set if it is successful
func (db *Database) CreateStudent(student *Student) {
	result, err := db.Collections.Students.InsertOne(context.TODO(), student)
	if err != nil {
		panic(err)
	}

	student.ID = result.InsertedID.(primitive.ObjectID)
}

// GetStudents returns a list of all students stored in the database.
func (db *Database) GetStudents() []Student {
	cur, err := db.Collections.Students.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	students := make([]Student, 0)
	if err := cur.All(context.TODO(), &students); err != nil {
		panic(err)
	}

	return students
}

// GetStudentByID gets a student by their ID. If not found, found will be false and
// err will be nil.
func (db *Database) GetStudentByID(id primitive.ObjectID) (student Student, found bool) {
	result := db.Collections.Students.FindOne(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Student{}, false
		}
		panic(err)
	}

	if err := result.Decode(&student); err != nil {
		panic(err)
	}

	found = true
	return
}

// DeleteStudent deletes a student from the database by ID. If the student could not be
// found, success will be false
func (db *Database) DeleteStudent(id primitive.ObjectID) bool {
	result := db.Collections.Students.FindOneAndDelete(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}
	return true
}

// UpdateStudent finds a student by its ID and updates it
func (db *Database) UpdateStudent(id primitive.ObjectID, newStudent *Student) bool {
	result := db.Collections.Locations.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{"$set": newStudent})
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}

	err = result.Decode(newStudent)

	return true
}

// GetStudentByHandle gets a student by the StudentHandles member. If the
// student is found, found will be true. If there is an error getting the student,
// it will be returned. Not that this error will not contain the not found error
func (db *Database) GetStudentByHandle(handle string) (student Student, found bool) {
	result := db.Collections.Students.FindOne(context.TODO(), bson.M{"studenthandles": bson.M{"$elemMatch": bson.M{"$eq": handle}}})

	err := result.Err()
	if err != nil {
		// If the student cannot be found
		if err == mongo.ErrNoDocuments {
			return Student{}, false
		}
		panic(err)
	}

	if err := result.Decode(&student); err != nil {
		panic(err)
	}

	found = true

	return
}

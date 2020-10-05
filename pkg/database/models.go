// +build generate

// This file contains generic code for implementing basic methods
// for each model such as references, Get by ID, Update, etc...
// If you're not modifying these functions, you shouldn't have to worry
// about regenerating code. However, if you updated this file, to update
// the changes for each of the models you would have to install
// genny https://github.com/cheekybits/genny and run go generate.

package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cheekybits/genny/generic"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//go:generate genny -in=$GOFILE -out=gen-student.go		-tag=generate gen "Model=Student model=student"
//go:generate genny -in=$GOFILE -out=gen-event.go		-tag=generate gen "Model=Event model=event"
//go:generate genny -in=$GOFILE -out=gen-location.go 	-tag=generate gen "Model=Location model=location"

type Model generic.Type

// ModelRef is a reference to a Model which, when serialized, will return
// the json of the referenced object.
//
// Be careful for circular references.
type ModelRef primitive.ObjectID

func (ref ModelRef) GetBSON() (interface{}, error) {
	return primitive.ObjectID(ref), nil
}

func (ref ModelRef) MarshalJSON() ([]byte, error) {
	obj, found := DB.GetModelByID(primitive.ObjectID(ref))
	if !found {
		return nil, fmt.Errorf("could not find Model with id %s", ref)
	}

	return json.Marshal(obj)
}

// Same functionality is ObjectID.UnmarshalJSON except it returns an error if the referenced object doesn't exist
func (ref *ModelRef) UnmarshalJSON(b []byte) error {
	id := primitive.ObjectID(*ref)
	if err := id.UnmarshalJSON(b); err != nil {
		return err
	}
	_, found := DB.GetModelByID(id)
	if !found {
		return fmt.Errorf("object not found")
	}

	*ref = ModelRef(id)
	return nil
}

// Gets the referenced object and panics if it doesn't exist
func (ref ModelRef) Get() Model {
	obj, found := DB.GetModelByID(primitive.ObjectID(ref))
	if !found {
		panic("could not find object")
	}
	return obj
}

// Ref creates a reference to the object
func (obj Model) Ref() ModelRef {
	return ModelRef(obj.ID)
}

// CreateModel creates a Model and adds it to the database. The
// ID element of the newly created Model will be set if it is successful
func (db *Database) CreateModel(model *Model) {
	result, err := db.Collections.Models.InsertOne(context.TODO(), model)
	if err != nil {
		panic(err)
	}

	model.ID = result.InsertedID.(primitive.ObjectID)
}

// GetModels returns a list of all models stored in the database.
func (db *Database) GetModels() []Model {
	cur, err := db.Collections.Models.Find(context.TODO(), bson.D{})
	if err != nil {
		panic(err)
	}

	models := make([]Model, 0)
	if err := cur.All(context.TODO(), &models); err != nil {
		panic(err)
	}

	return models
}

// GetModelByID gets a model by their ID. If not found, found will be false and
// err will be nil.
func (db *Database) GetModelByID(id primitive.ObjectID) (model Model, found bool) {
	result := db.Collections.Models.FindOne(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Model{}, false
		}
		panic(err)
	}

	if err := result.Decode(&model); err != nil {
		panic(err)
	}

	found = true
	return
}

// GetModelByIDString gets a model by its ID as a string. If the ID could not be
// parsed into an object ID or the model could not be found, an error will be returned
func (db *Database) GetModelByIDString(id string) (Model, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Model{}, err
	}

	model, found := db.GetModelByID(objectID)
	if !found {
		return Model{}, fmt.Errorf("no models found with id %s", objectID.Hex())
	}

	return model, nil
}

// DeleteModel deletes a model from the database by ID. If the model could not be
// found, success will be false
func (db *Database) DeleteModel(id primitive.ObjectID) bool {
	result := db.Collections.Models.FindOneAndDelete(nil, bson.M{"_id": id})

	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}
	return true
}

// UpdateModel finds a model by its ID and updates it
func (db *Database) UpdateModel(id primitive.ObjectID, newModel *Model) bool {
	result := db.Collections.Locations.FindOneAndUpdate(nil, bson.M{"_id": id}, bson.M{"$set": newModel})
	err := result.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false
		}
		panic(err)
	}

	err = result.Decode(newModel)

	return true
}

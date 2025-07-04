package mongodb

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// Collections holds references to initialized MongoDB collections
var Collections = make(map[string]*mongo.Collection)

// CollectionNames defines the names of all collections
var CollectionNames = map[string]string{
	"devices": "devices",
}

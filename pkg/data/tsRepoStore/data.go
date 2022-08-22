package tsRepoStore

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// UserData Setting up the database for the user data collection
func UserData(dbClient *mongo.Client, collectionName string) *mongo.Collection {
	var userCollection = dbClient.Database("track_space").Collection(collectionName)
	return userCollection
}

// MailData Setting up the database for the mail data collection
func MailData(dbClient *mongo.Client, collectionName string) *mongo.Collection {
	var mailCollection = dbClient.Database("track_space").Collection(collectionName)
	return mailCollection
}

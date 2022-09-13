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
func AdminData(dbClient *mongo.Client, collectionName string) *mongo.Collection {
	var Admin = dbClient.Database("track_space").Collection(collectionName)
	return Admin
}

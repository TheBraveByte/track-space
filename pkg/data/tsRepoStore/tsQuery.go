package tsRepoStore

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (tm TsMongoDBRepo) InsertInfo(email, password string) (int64, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: email}}
	count, err := UserData(tm.TsMongoDB, "user").CountDocuments(ctx, filter)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if count == 1 {
		return 0, nil
	}

	if count == 0 {
		documents := bson.D{
			{Key: "email", Value: email},
			{Key: "password", Value: password},
		}
		_, err = UserData(tm.TsMongoDB, "user").InsertOne(ctx, documents)

		if err != nil {
			log.Println("cannot insert user sign up details in the database")
			return 0, err
		}
	}
	return 0, nil
}

func (tm TsMongoDBRepo) UpdateUserInfo(info map[string]interface{}, email interface{}, t1, t2 string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: email}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "first_name", Value: info["firstname"]},
		{Key: "last_name", Value: info["lastname"]},
		{Key: "address", Value: info["address"]},
		{Key: "yrs_of_exp", Value: info["YrsOfExp"]},
		{Key: "country", Value: info["country"]},
		{Key: "stack", Value: info["stack"]},
		{Key: "phone_number", Value: info["phone"]},
		{Key: "ip_address", Value: info["IPAddress"]},
		{Key: "created_at", Value: info["created_at"]},
		{Key: "updated_at", Value: info["updated_at"]},
		{Key: "token", Value: t1},
		{Key: "renew_token", Value: t2},
	}}}
	var updateDocument bson.M
	err := UserData(tm.TsMongoDB, "user").FindOneAndUpdate(ctx, filter, update).Decode(&updateDocument)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("cannot find and update the database", err)
			return err
		}
		log.Fatal(err)
	}
	return nil

}

func (tm TsMongoDBRepo) VerifyLogin(email string) (bool, string) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()

	var result bson.M
	filter := bson.D{{Key: "email", Value: email}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "No previous user document match found"
		}
		log.Fatal(err)
		return false, "No match document found"
	}
	password := fmt.Sprintf("%v", result["password"])
	return true, password
}

func (tm TsMongoDBRepo) SendUserDetails(email interface{})(primitive.M, error){
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var result bson.M
	filter := bson.D{{Key: "email", Value: email}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		log.Fatal("cannot find document")
		return nil, err
	}
	
	return result, nil
}




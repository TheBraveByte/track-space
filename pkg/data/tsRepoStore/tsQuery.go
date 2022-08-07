package tsRepoStore

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (tm *TsMongoDBRepo) InsertInfo(email, password string) (int64, error) {
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

func (tm *TsMongoDBRepo) UpdateUserInfo(user model.User, email interface{}, t1, t2 string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()

	filter := bson.D{{Key: "email", Value: email}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "first_name", Value: user.FirstName },
		{Key: "last_name", Value:user.LastName },
		{Key: "address", Value: user.Address},
		{Key: "yrs_of_exp", Value:user.YrsOfExp},
		{Key: "country", Value: user.Country},
		{Key: "stack", Value: user.Stack},
		{Key: "phone_number", Value: user.PhoneNumber},
		{Key: "ip_address", Value: user.IPAddress},
		{Key: "created_at", Value: user.CreatedAt},
		{Key: "updated_at", Value: user.UpdatedAt},
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

func (tm *TsMongoDBRepo) VerifyLogin(email string) (bool, string) {
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

func (tm *TsMongoDBRepo) SendUserDetails(email interface{}) (primitive.M, error) {
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

func (tm *TsMongoDBRepo) StoreWorkSpaceData(email interface{}, project model.Project) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	var EndTime time.Time
	project.EndTime =  EndTime.Local().UTC()
	var result bson.M
	filter := bson.D{{Key: "email", Value: email}}
	update := bson.D{
		{Key: "project_details", Value: bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "project_name", Value: project.ProjectName},
				{Key: "tools_use_as", Value: project.ToolsUseAs},
				{Key: "project_content", Value: project.ProjectContent},
				{Key: "start_time", Value: project.StartTime},
				{Key: "end_time", Value: project.EndTime},
			}},
		}},
	}
	err := UserData(tm.TsMongoDB, "user").FindOneAndUpdate(ctx, filter, update).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return err
		}
		log.Fatal("cannot find document")
		return err
	}
	return nil
}

func (tm *TsMongoDBRepo) OrganizeWorkSpaceData(projectData model.User, email string)(map[string]int, error){
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	var count map[string]int

	var result bson.M
	filter := bson.D{{Key: "email", Value: email}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return count, err
		}
		log.Fatal("cannot find document")
		return count, err
	}

	for key , value := range result{
		if key == "project_details" {
			switch v := value.(type){
			
			case []model.Project:
				var countCode, countText int = 0, 0
				for _, y :=range v{
					if y.ToolsUseAs == "code"{
						countCode +=1
						
					} else {
						countText += 1
					}
				}
				newFunction(count, countCode)
				newFunction1(count, countText)
			}
					
		}
	}
	return count, nil
}

func newFunction1(count map[string]int, countText int) {
	count["text"] = countText
}

func newFunction(count map[string]int, countCode int) {
	count["code"] = countCode
}
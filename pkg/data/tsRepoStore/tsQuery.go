package tsRepoStore

import (
	"context"
	"log"
	"time"

	"github.com/yusuf/track-space/pkg/key"
	"github.com/yusuf/track-space/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (tm *TsMongoDBRepo) InsertUserInfo(user_id, email, password string) (int64, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var userInfo bson.M
	filter := bson.D{
		{Key: "_id", Value: user_id},
	}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&userInfo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			documents := bson.D{
				{Key: "_id", Value: user_id},
				{Key: "email", Value: email},
				{Key: "password", Value: password},
			}
			_, err := UserData(tm.TsMongoDB, "user").InsertOne(ctx, documents)
			if err != nil {
				panic(err)
			}
			return 0, nil
		}
		panic(err)
	}
	return 1, nil
}

func (tm *TsMongoDBRepo) UpdateUserInfo(user model.User, id, t1, t2 string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: id}}

	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "first_name", Value: user.FirstName},
		{Key: "last_name", Value: user.LastName},
		{Key: "address", Value: user.Address},
		{Key: "yrs_of_exp", Value: user.YrsOfExp},
		{Key: "country", Value: user.Country},
		{Key: "stack", Value: user.Stack},
		{Key: "phone_number", Value: user.PhoneNumber},
		{Key: "ip_address", Value: user.IPAddress},
		{Key: "created_at", Value: user.CreatedAt},
		{Key: "updated_at", Value: user.UpdatedAt},
		{Key: "token", Value: t1},
		{Key: "renew_token", Value: t2},
	}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("cannot find and update the database", err)
			return err
		}
		log.Fatal(err)
	}
	return nil
}

func (tm *TsMongoDBRepo) UpdateUserField(id, t1, t2 string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "token", Value: t1}, 
		{Key: "renew_token", Value: t2},
	}}}
	// var updateDocument bson.M
	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println("cannot find and update the database", err)
			return err
		}
		log.Fatal(err)
	}
	return nil
}

func (tm *TsMongoDBRepo) VerifyLogin(id, hashedPassword, postPassword string) (bool, string) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)

	defer cancelCtx()

	var result bson.M
	filter := bson.D{
		{Key: "_id", Value: id},
	}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "not a registered user"
		}
		log.Fatal(err)
		return false, "No match document found"
	}
	ok, msg := key.VerifyPassword(postPassword, hashedPassword)
	return ok, msg
}

func (tm *TsMongoDBRepo) SendUserDetails(id string) (primitive.M, error) {
	// this was use twice in the controllers
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var user bson.M
	filter := bson.D{{Key: "_id", Value: id}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		log.Panic("cannot find document")

	}

	return user, nil
}

func (tm *TsMongoDBRepo) StoreWorkSpaceData(id string, project model.Project) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "project_details", Value: bson.D{
			{Key: "_id", Value: project.ID},
			{Key: "project_name", Value: project.ProjectName},
			{Key: "tools_use_as", Value: project.ToolsUseAs},
			{Key: "project_content", Value: project.ProjectContent},
			{Key: "created_at", Value: project.CreatedAt},
			{Key: "updated_at", Value: project.UpdatedAt},
			{Key: "status", Value: project.Status},
		}},
	}}}
	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("cannot find document")
		return err
	}
	return nil
}

func (tm *TsMongoDBRepo) ModifyProjectData(id string, project model.Project) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "project_details", Value: bson.D{
			{Key: "project_name", Value: project.ProjectName},
			{Key: "tools_use_as", Value: project.ToolsUseAs},
			{Key: "project_content", Value: project.ProjectContent},
			{Key: "created_at", Value: project.CreatedAt},
			{Key: "updated_at", Value: project.UpdatedAt},
		}},
	}}}
	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("cannot find document")
		return err
	}
	return nil
}

func (tm *TsMongoDBRepo) StoreDailyTaskData(task model.DailyTask, id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "todo", Value: bson.D{
			{Key: "_id", Value: task.ID},
			{Key: "to_do_task", Value: task.ToDoTask},
			{Key: "date_schedule", Value: task.DateSchedule},
			{Key: "start_time", Value: task.StartTime},
			{Key: "end_time", Value: task.EndTime},
		}},
	}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("cannot add document to the database")
		return err
	}
	return nil
}

func (tm *TsMongoDBRepo) GetProjectData(project_id string) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	
	filter := bson.D{{"project_details.project_name", project_name}}
	opt :=  options.Find().SetProjection(bson.D{
		{"_id", "63134643e3c4d8b40db322f2"},
		{"project_details", bson.D{{"$elemMatch", bson.D{{"project_name", "newcode"}}}}},
	})
	}
	// opt :=  options.Find().SetProjection() //bson.M{"project_details":bson.M{"$elemMatch": bson.M{"_id":project_id}}}
	var data bson.M
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&data)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			log.Panic(err)
		}
	}
	log.Println(len(data))
	log.Println(data)
	// for k,v :=range data{
		
	// }
	return data, nil
}

func (tm *TsMongoDBRepo) UpdateProjectStat(data model.Data, id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "data", Value: bson.D{
		{Key: "date", Value: data.Date},
		{Key: "code", Value: data.Code},
		{Key: "article", Value: data.Article},
		{Key: "text", Value: data.Text},
		{Key: "todo", Value: data.Todo},
		{Key: "total", Value: data.Total},
	}}}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Println("cannot update user document")
		panic(err)
	}
	return nil
}

func (tm *TsMongoDBRepo) GetProjectStatByID(id string) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var result bson.M
	filter := bson.D{{Key: "_id", Value: id}}
	err := UserData(tm.TsMongoDB, "stat").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err)
			panic(err)
		}
		return nil, err
	}

	return result, nil
}

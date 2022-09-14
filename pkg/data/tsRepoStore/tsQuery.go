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

/*
InsertUserInfo : this will help create a document for every user that sign up
on track space
*/
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
				log.Panicf("Error 0 from InsertUserInfo: %v", err)
			}
			return 0, nil
		}
		log.Printf("Error 1 from InsertUserInfo: %v", err)
	}
	return 1, nil
}

/*
UpdateUserInfo : this is to update a particular user document previous stored in the
database to add more information about the user
*/
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
		{Key: "profession", Value: user.Profession},
		{Key: "token", Value: t1},
		{Key: "renew_token", Value: t2},
	}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Fatalf("Error 0 from UpdateUserInfo: %v", err)
		}
		log.Fatalf("Error 1 from UpdateUserInfo: %v", err)
	}
	return nil
}

/*
UpdateUserField : this is to update the user generated token when signing in into track space
*/
func (tm *TsMongoDBRepo) UpdateUserField(id, t1, t2 string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "token", Value: t1}, {Key: "renew_token", Value: t2}}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error  from UpdateUserField: %v", err)
		return err
	}
	return nil
}

/*
VerifyLogin: this method will help to verify the user login input details with respect to
the store details in the database for authenication
*/
func (tm *TsMongoDBRepo) VerifyLogin(id, hashedPassword, postPassword string) (bool, string) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var result bson.M
	filter := bson.D{{Key: "_id", Value: id}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "no document found"
		}
		return false, "No match document found"
	}
	ok, msg := key.VerifyPassword(postPassword, hashedPassword)
	return ok, msg
}

/*
SendUserDetails : this method will help in getting a the user store information and
activities on trackspace when the user details is needed
*/
func (tm *TsMongoDBRepo) SendUserDetails(id string) (primitive.M, error) {
	// this was called  multiple time  in the controllers package
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var user bson.M
	filter := bson.D{{Key: "_id", Value: id}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		log.Fatalf("Error from SendUserDetails : %v", err)

	}
	return user, nil
}

/*
StoreProjectData : this method help the user to store the create project and all it
content on the workspace to the database
*/
func (tm *TsMongoDBRepo) StoreProjectData(id string, project model.Project) error {
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
		log.Fatalf("Error from StoreProjectData : %v", err)
	}
	return nil
}

/*
GetProjectData : this method fetch one particular created projects stored by a
particular user in the database to check or make some modification to the projects
*/
func (tm *TsMongoDBRepo) GetProjectData(project_id string) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "project_details._id", Value: project_id}}
	opt := options.FindOne().SetProjection(bson.D{
		{Key: "_id", Value: project_id},
		{Key: "project_details", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "_id", Value: project_id}}}}},
	})

	var data bson.M
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter, opt).Decode(&data)
	if err != nil {
		log.Fatalf("Error from GetProjectData: %v", err)
	}
	return data, nil
}

/*
ModifyProjectData : this method is to keep track of the changes made by the
user on a particular project by updating it in the database
*/
func (tm *TsMongoDBRepo) ModifyProjectData(user_id, id string, project model.Project) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{
		{Key: "_id", Value: user_id},
		{Key: "project_details._id", Value: id},
	}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "project_details.$._id", Value: project.ID},
		{Key: "project_details.$.project_name", Value: project.ProjectName},
		{Key: "project_details.$.tools_use_as", Value: project.ToolsUseAs},
		{Key: "project_details.$.project_content", Value: project.ProjectContent},
		{Key: "project_details.$.updated_at", Value: project.UpdatedAt},
		{Key: "project_details.$.status", Value: project.Status},
	}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Error from ModifyProjectData: %v", err)
	}
	return nil
}

/*
DeleteUserProject : this method will delete a select project by the user
*/
func (tm *TsMongoDBRepo) DeleteUserProject(project_id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "project_details._id", Value: project_id}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "project_details", Value: bson.D{{Key: "_id", Value: project_id}}},
	}}}
	opt := options.Update().SetUpsert(false)
	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update, opt)
	if err != nil {
		log.Fatalf("Error from DeleteUserProject : %v", err)
	}
	return nil
}

/*
StoreTodoData : this method help the user to store the create todo schedule and all it
set duration and date as well in the to the database
*/
func (tm *TsMongoDBRepo) StoreTodoData(todo model.Todo, id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$push", Value: bson.D{
		{Key: "todo", Value: bson.D{
			{Key: "_id", Value: todo.ID},
			{Key: "to_do_task", Value: todo.ToDoTask},
			{Key: "date_schedule", Value: todo.DateSchedule},
			{Key: "start_time", Value: todo.StartTime},
			{Key: "end_time", Value: todo.EndTime},
		}},
	}}}

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatal("cannot add document to the database")
		return err
	}
	return nil
}

/*
GetTodoData : this method fetch one particular created schedule stored by a
particular user in the database to check or make some modification to the
date/time of the schedule
*/
func (tm *TsMongoDBRepo) GetTodoData(todo_id string) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "todo._id", Value: todo_id}}
	opt := options.FindOne().SetProjection(bson.D{
		{Key: "_id", Value: todo_id},
		{Key: "todo", Value: bson.D{{Key: "$elemMatch", Value: bson.D{{Key: "_id", Value: todo_id}}}}},
	})

	var data bson.M
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter, opt).Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	return data, nil
}

/*
ModifyTodoData : this method is to keep track of the changes made by the
user on a previous set schedule by updating it in the database
*/
func (tm *TsMongoDBRepo) ModifyTodoData(id string, todo model.Todo) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()
	filter := bson.D{{Key: "todo._id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "todo.$._id", Value: todo.ID},
		{Key: "todo.$.to_do_task", Value: todo.ToDoTask},
		{Key: "todo.$.date_schedule", Value: todo.DateSchedule},
		{Key: "todo.$.start_time", Value: todo.StartTime},
		{Key: "todo.$.end_time", Value: todo.EndTime},
	}}}
	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf("Error from ModifyTodoData : %v", err)
	}
	return nil
}

/*
DeleteUserProject : this method will delete a select todo schedule by the user
*/
func (tm *TsMongoDBRepo) DeleteUserTodo(todo_id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "todo._id", Value: todo_id}}
	update := bson.D{{Key: "$pull", Value: bson.D{
		{Key: "todo", Value: bson.D{{Key: "_id", Value: todo_id}}},
	}}}
	opt := options.Update().SetUpsert(false)

	_, err := UserData(tm.TsMongoDB, "user").UpdateOne(ctx, filter, update, opt)
	if err != nil {
		log.Fatalf("Error from DeleteUserTodo : %v", err)
	}
	return nil
}

/*
UpdateOne : this method is to store the statistic updates of the user activities
on track space
*/
func (tm *TsMongoDBRepo) UpdateUserStat(data model.Data, id string) error {
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

/*
GetUserStatByID : this method is to get all the statistic information on a
particular user
*/
func (tm *TsMongoDBRepo) GetUserStatByID(id string) (primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var result bson.M
	filter := bson.D{{Key: "_id", Value: id}}
	err := UserData(tm.TsMongoDB, "user").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Println(err)
		}
	}

	return result, nil
}

func (tm *TsMongoDBRepo) GetAllUserData() ([]primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var documents []bson.M
	cursor, err := UserData(tm.TsMongoDB, "user").Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}
	if err = cursor.All(ctx, &documents); err != nil {
		log.Panic(err)
		return nil, err
	}

	return documents, nil
}

func (tm *TsMongoDBRepo) GetAdminInfo() ([]primitive.M, error) {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	var result []bson.M
	cursor, err := AdminData(tm.TsMongoDB, "admin").Find(ctx, bson.D{})
	if err != nil {
		log.Panic(err)
	}

	if err = cursor.All(ctx, &result); err != nil {
		log.Panic(err)
		return nil, err
	}
	return result, nil
}

func (tm *TsMongoDBRepo) AdminDeleteUserData(id string) error {
	ctx, cancelCtx := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancelCtx()

	filter := bson.D{{Key: "_id", Value: id}}
	var deletedProject bson.M
	err := UserData(tm.TsMongoDB, "user").FindOneAndDelete(ctx, filter).Decode(&deletedProject)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Panic(err)
			return err
		}
		log.Fatal(err)
		return err
	}
	return nil
}

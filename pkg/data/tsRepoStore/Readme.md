# Database Query Documentation
This is a documentation for the database query package tsRepoStore. This package provides a set of methods that interact with MongoDB to manipulate data related to users of the **track-space** application.

### Installation
To install the package, use the following command:

`go get github repository link`

### Usage
To use the package, you need to import it into your code as follows:

You also need to provide a valid MongoDB instance and database to the package. This can be done using the TsMongoDBRepo struct. Here's an example:

```go

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		panic(err)
	}

	db := client.Database("mydb")
	repo := tsRepoStore.TsMongoDBRepo{TsMongoDB: db}
	// Use repo methods here...
}
```
### Methods

#### InsertUserInfo
`go
func (tm *TsMongoDBRepo) InsertUserInfo(email, password string) (int64, string, error)
`

This method creates a new document in the user collection to represent a new user. The document contains the user's email and password. It returns an int64 that represents the number of documents created and a string that represents the ID of the created document. An error is returned if there is an issue creating the document.

#### UpdateUserInfo
`go
func (tm *TsMongoDBRepo) UpdateUserInfo(user model.User, id, t1, t2 string) error
`

This method updates an existing user's document with additional user information. The method takes a model.User object that represents the new user information, the id of the user's document, and two tokens t1 and t2. The method updates the user's document with the new information and tokens. An error is returned if there is an issue updating the document.

#### UpdateUserField
`go
func (tm *TsMongoDBRepo) UpdateUserField(id, t1, t2 string) error
`

This method updates an existing user's document with new tokens. The method takes the id of the user's document and two tokens t1 and t2. The method updates the user's document with the new tokens. An error is returned if there is an issue updating the document.

#### ResetUserPassword
`go
func (tm *TsMongoDBRepo) ResetUserPassword(email, newPassword string) error
`

This method resets an existing user's password. The method takes the email of the user and the new newPassword. The method updates the user's document with the new password. An error is returned if there is an issue updating the document.

#### VerifyLogin
`go
func (tm *TsMongoDBRepo) VerifyLogin(email, password string) (bool, string, string, error)
`

This method verifies a user's login credentials. The method takes the user's email and password. The method returns a boolean indicating whether the credentials are correct, the id of the user's document, and two tokens `t1




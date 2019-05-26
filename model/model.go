package model

import (
	"context"
	"log"
	"time"

	"github.com/mitchellh/mapstructure"
	"golang.org/x/crypto/bcrypt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

// UserInfo : UserInfo-class
type UserInfo struct {
	FirstName string
	LastName  string
	CreatedAt time.Time
}

// User : User-class
type User struct {
	UserName     string
	PasswordHash string
	Info         UserInfo
}

// UserFactory : In charge of producing user objects
type UserFactory struct{}

// CreateUser : Crest
func (uf UserFactory) CreateUser(UserName string, FirstName string, LastName string, Password string, CreatedAt time.Time) User {
	info := UserInfo{FirstName, LastName, CreatedAt}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(Password), 4)
	if err != nil {
		log.Fatalf("Can't hash password: %v\n", err)
	}

	passwordHashString := string(passwordHash)
	user := User{UserName, passwordHashString, info}
	return user
}

// UserRepository : In charge of communicating with user-db
type UserRepository struct{}

// WriteUserToDb : Stores user-object in database
func (ur UserRepository) WriteUserToDb(user User) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	_, err = client.Collection("Users").Doc(user.UserName).Set(ctx, user)
	if err != nil {
		log.Fatalf("error writing to firebase: %v\n", err)
	}
}

// GetUserFromDb : Retrieves user object from database by username
func (ur UserRepository) GetUserFromDb(UserName string) User {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	var user User
	snapshot, err := client.Collection("Users").Doc(UserName).Get(ctx)
	if err != nil {
		log.Fatalf("error retrieving user from database: %v\n", err)
	}
	mapstructure.Decode(snapshot.Data(), &user)

	return user
}

// IsUserNameTaken : Returns true if taken, false if available
func (ur UserRepository) IsUserNameTaken(UserName string) bool {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	_, err = client.Collection("Users").Doc(UserName).Get(ctx)
	if err != nil {
		return false
	}

	return true
}

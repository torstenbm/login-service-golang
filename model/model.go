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
func (ur UserRepository) WriteUserToDb(user User) error {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return err
	}

	_, err = client.Collection("Users").Doc(user.UserName).Set(ctx, user)
	if err != nil {
		return err
	}
	return err
}

// GetUserFromDb : Retrieves user object from database by username
func (ur UserRepository) GetUserFromDb(UserName string) (User, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")
	var user User

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return user, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return user, err
	}

	snapshot, err := client.Collection("Users").Doc(UserName).Get(ctx)
	if err != nil {
		return user, err
	}
	mapstructure.Decode(snapshot.Data(), &user)

	return user, err
}

// IsUserNameTaken : Returns true if taken, false if available
func (ur UserRepository) IsUserNameTaken(UserName string) (bool, error) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("model/go-login-service-firebase-adminsdk-rkpp7-075307324f.json")

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return true, err
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		return true, err
	}

	_, err = client.Collection("Users").Doc(UserName).Get(ctx)
	if err != nil {
		return false, err
	}

	return true, err
}

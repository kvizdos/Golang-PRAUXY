package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var conn = ConnectDB("users")

var rdb = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

var ctx = context.Background()

type User struct {
	Id           string         `bson:"id,omitempty"`
	Username     string         `bson:"username,omitempty"`
	Email        string         `bson:"email,omitempty"`
	Password     string         `bson:"password,omitempty"`
	Multifactor  []mfaInterface `bson:"multifactor,omitempty"`
	SessionToken string         `bson:"session_token,omitempty"`
}

type CRUD interface {
	Create(User) (struct{}, error)
	Read(User) User
	Update(User) error
	Delete(User) error
}

type UserCRUD struct {
	// Create(User) error
	// Read(page, size, skip) []User
	// Update(User) error
	// Delete(User) error
}

type Response struct {
	statusCode int
	body       string
}

// https://gowebexamples.com/password-hashing/
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (uc *UserCRUD) Create(u User) (*mongo.InsertOneResult, error) {
	stat, err := conn.InsertOne(context.TODO(), u)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return stat, nil
}

func (uc *UserCRUD) Read(u User) (User, error) {
	var result User

	err := conn.FindOne(context.TODO(), u).Decode(&result)

	return result, err
}

func RegisterUser(user User) Response {
	crud := UserCRUD{}

	u := User{
		Username: strings.ToLower(user.Username),
	}

	_, err := crud.Read(u)

	if err != mongo.ErrNoDocuments {
		r := Response{
			statusCode: 409,
			body:       "username taken",
		}
		return r
	}

	hashedPass, _ := hashPassword(user.Password)
	user.Password = hashedPass
	user.Username = strings.ToLower(user.Username)

	crud.Create(user)

	r := Response{
		statusCode: 200,
		body:       fmt.Sprintf("user %s created", user.Username),
	}

	return r
}
func RegisterUserHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	u := User{
		Username: body["username"],
		Email:    body["email"],
		Password: body["password"],
	}

	status := RegisterUser(u)

	w.WriteHeader(status.statusCode)
	fmt.Fprint(w, status.body)
}

// https://stackoverflow.com/questions/45267125/how-to-generate-unique-random-alphanumeric-tokens-in-golang
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func updateUserSessionToken(user User) string {
	newToken := GenerateSecureToken(128)

	rawToken := []byte(newToken)
	hashedToken := sha256.Sum256(rawToken)

	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": bson.M{"session_token": fmt.Sprintf("%x", hashedToken)}}

	conn.UpdateOne(context.TODO(), filter, update)
	return newToken
}

func LoginUser(validateUser User) Response {
	crud := UserCRUD{}

	checkForUsername := User{
		Username: validateUser.Username,
	}
	foundUser, err := crud.Read(checkForUsername)

	if err == mongo.ErrNoDocuments {
		r := Response{
			statusCode: 401,
			body:       "invalid username or password",
		}
		return r
	}

	if checkPasswordHash(validateUser.Password, foundUser.Password) {
		r := Response{}
		token := updateUserSessionToken(foundUser)

		if len(foundUser.Multifactor) == 0 {
			r.statusCode = 200
			r.body = fmt.Sprintf(`{"token": "%s"}`, token)
		} else {
			mfaSID := GenerateSecureToken(32)
			rdb.Set(ctx, fmt.Sprintf("%x", sha256.Sum256([]byte(mfaSID))), token, 0)
			r.statusCode = 200
			r.body = fmt.Sprintf(`{"mfaSID": "%s"}`, mfaSID)
		}
		return r
	}

	r := Response{
		statusCode: 401,
		body:       "invalid username or password",
	}
	return r
}
func LoginUserHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	user := User{
		Username: strings.ToLower(body["username"]),
		Password: body["password"],
	}

	status := LoginUser(user)

	w.WriteHeader(status.statusCode)
	fmt.Fprint(w, status.body)
}

func AuthorizeUser(username string, token string) {
	return
}

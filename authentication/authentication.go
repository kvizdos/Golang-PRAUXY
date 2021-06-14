package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	// "github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var conn = ConnectDB("users")

type User struct {
	Username    string         `bson:"username,omitempty"`
	Email       string         `bson:"email,omitempty"`
	Password    string         `bson:"password,omitempty"`
	Multifactor []mfaInterface `bson:"multifactor,omitempty"`
}

type mfaInterface struct {
	typeOfMfa string // Available: totp, fido
	secret    string
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
		Username: user.Username,
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

	// totpOpts := totp.GenerateOpts{
	// 	Issuer:      "PRAUXY",
	// 	AccountName: user.Username,
	// }
	// key, _ := totp.Generate(totpOpts)

	user.Password = hashedPass
	// user.TotpSecret = key.Secret()

	crud.Create(user)

	// var buf bytes.Buffer
	// img, _ := key.Image(200, 200)
	// png.Encode(&buf, img)

	r := Response{
		statusCode: 200,
		body:       "session token somewhere here",
		// body:       fmt.Sprintf(`{"totp": "%s"}`, b64.StdEncoding.EncodeToString(buf.Bytes())),
	}

	return r
}
func RegisterUserHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	u := User{
		Username: strings.ToLower(body["username"]),
		Email:    body["email"],
		Password: body["password"],
	}

	status := RegisterUser(u)

	w.WriteHeader(status.statusCode)
	fmt.Fprint(w, status.body)
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
			body:       "invalid username",
		}
		return r
	}

	if checkPasswordHash(validateUser.Password, foundUser.Password) {
		r := Response{
			statusCode: 200,
			body:       "session token somewhere here",
		}

		return r
	}

	r := Response{
		statusCode: 401,
		body:       "invalid password",
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

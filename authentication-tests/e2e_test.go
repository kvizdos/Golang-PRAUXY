package main

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	Id          string         `bson:"id,omitempty"`
	Username    string         `bson:"username,omitempty"`
	Email       string         `bson:"email,omitempty"`
	Password    string         `bson:"password,omitempty"`
	Multifactor []mfaInterface `bson:"multifactor,omitempty"`
}

type mfaInterface struct {
	TypeOfMFA string `bson:"type"` // Available: totp, fido
	Secret    string `bson:"secret"`
}

func ConnectDB(collectionName string) *mongo.Collection {
	clientOpts := options.Client().ApplyURI("mongodb://auth_mongo")
	client, err := mongo.Connect(context.TODO(), clientOpts)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("prauxy").Collection(collectionName)

	return collection
}

var testConn = ConnectDB("users")

func makeReq(endpoint string, method string, json string, t *testing.T) (int, string) {
	url := fmt.Sprintf("http://backend:8080%s", endpoint)

	var jsonStr = []byte(json)
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return resp.StatusCode, string(body)
}

// Registration Tests
func TestRegistrationEndpointFailsWithBadBody(t *testing.T) {
	statusCode, body := makeReq("/register", "POST", `{}`, t)

	if statusCode != 400 || strings.Compare(body, "Missing body field username\n") != 0 {
		t.Fatalf("Received code %s (expected 400) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestRegistrationNewUser(t *testing.T) {
	statusCode, body := makeReq("/register", "POST", `{
		"username": "kenton",
		"email": "kvizdos@gmail.com",
		"password": "password123"
	}`, t)

	if statusCode != 200 || strings.Compare(body, "user kenton created") != 0 {
		t.Fatalf("Received code %s (expected 200) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestRegistrationExistingUser(t *testing.T) {
	statusCode, body := makeReq("/register", "POST", `{
		"username": "kenton",
		"email": "kvizdos@gmail.com",
		"password": "password123"
	}`, t)

	if statusCode != 409 || strings.Compare(body, "username taken") != 0 {
		t.Fatalf("Received code %s (expected 409) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestConfirmRegistrationHashesPassword(t *testing.T) {
	var result User

	testConn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	if result.Password == "password123" {
		t.Fatalf("Password is not hashed!")
	}
}

func TestRegistrationExistingUserWithDifferentCapitalization(t *testing.T) {
	statusCode, body := makeReq("/register", "POST", `{
		"username": "KenTon",
		"email": "kvizdos@gmail.com",
		"password": "password123"
	}`, t)

	if statusCode != 409 || strings.Compare(body, "username taken") != 0 {
		t.Fatalf("Received code %s (expected 409) with error: %s", strconv.Itoa(statusCode), body)
	}
}

// Login Tests
func TestLoginWithInvalidUsername(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "this_person_does_not_exist",
		"password": "password123"
	}`, t)

	if statusCode != 401 || strings.Compare(body, "invalid username or password") != 0 {
		t.Fatalf("Received code %s (expected 401) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidUsernameButBadPassword(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "bad_pass"
	}`, t)

	if statusCode != 401 || strings.Compare(body, "invalid username or password") != 0 {
		t.Fatalf("Received code %s (expected 401) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidUsernameAndPassword(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "password123"
	}`, t)

	match, _ := regexp.MatchString(`{"token": ".+"}`, body)

	if statusCode != 200 || match == false {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidPasswordButOddlyCapitalizedUsername(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "KenToN",
		"password": "password123"
	}`, t)

	match, _ := regexp.MatchString(`{"token": ".+"}`, body)

	if statusCode != 200 || match == false {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

// TOTP Tests
func TestTOTPFailsWhenUsingInvalidType(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "POST", `{
		"type": "bad_type",
		"username": "kenton"
	}`, t) // Normally available: totp, fido

	if statusCode != 406 || strings.Compare(body, "invalid mfa type") != 0 {
		t.Fatalf("Received code %s (expected 406) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPCreationSuccessful(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "POST", `{
		"type": "totp",
		"username": "kenton"
	}`, t)

	match, _ := regexp.MatchString(`{"qr": ".+"}`, body)

	var result User

	testConn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	if statusCode != 200 || match == false || (len(result.Multifactor) == 0 || result.Multifactor[0].TypeOfMFA != "totp") {
		t.Fatalf("Received code %s (expected 200). Len of mfa: %s, mfaInfo: %v", strconv.Itoa(statusCode), strconv.Itoa(len(result.Multifactor)), result.Multifactor[0].TypeOfMFA)
	}
}

func TestTOTPFailsWhenTryingToCreateSecondOTP(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "POST", `{
		"type": "totp",
		"username": "kenton"
	}`, t)

	if statusCode != 406 || strings.Compare(body, "totp already registered") != 0 {
		t.Fatalf("Received code %s (expected 406) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPDeleteFailsWithInvalidType(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "DELETE", `{
		"type": "bad_type",
		"username": "kenton"
	}`, t)

	if statusCode != 406 || strings.Compare(body, "invalid mfa type") != 0 {
		t.Fatalf("Received code %s (expected 406) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPDeleteSucceeds(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "DELETE", `{
		"type": "totp",
		"username": "kenton"
	}`, t)

	var result User

	testConn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	if statusCode != 200 || strings.Compare(body, "totp disabled") != 0 || len(result.Multifactor) != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s and multifactor length of %d", strconv.Itoa(statusCode), body, len(result.Multifactor))
	}
}

func TestTOTPDeleteFailsWithDisabledType(t *testing.T) {
	statusCode, body := makeReq("/user/mfa", "DELETE", `{
		"type": "totp",
		"username": "kenton"
	}`, t)

	if statusCode != 406 || strings.Compare(body, "totp not enabled") != 0 {
		t.Fatalf("Received code %s (expected 406) with body: %s", strconv.Itoa(statusCode), body)
	}
}

// MFA Login Tests
func TestTOTPValidationFailsWithBadSID(t *testing.T) {
	makeReq("/user/mfa", "POST", `{
		"type": "totp",
		"username": "kenton"
	}`, t)

	makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "password123"
	}`, t)

	statusCode, body := makeReq("/user/mfa/verify", "POST", fmt.Sprintf(`{
		"type": "totp",
		"username": "%s",
		"sid": "dfgijdfgi",
		"code": "doesn't matter yet"
	}`, "kenton"), t)

	if statusCode != 403 || strings.Compare(body, "invalid mfa sid") != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPValidationFailsWithBadCodeButGoodSID(t *testing.T) {
	var result User

	testConn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	ti := time.Now()
	passcode, _ := totp.GenerateCode(result.Multifactor[0].Secret, ti)

	_, bodyLogin := makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "password123"
	}`, t)

	mfaSIDRgx := regexp.MustCompile(`{"mfaSID": "(.*?)"}`)
	SID := mfaSIDRgx.FindStringSubmatch(bodyLogin)

	statusCode, body := makeReq("/user/mfa/verify", "POST", fmt.Sprintf(`{
		"type": "totp",
		"username": "kenton",
		"sid": "%s",
		"code": "%s"
	}`, SID[1], passcode), t)

	match, _ := regexp.MatchString(`{"token": ".+"}`, body)

	if statusCode != 200 || match == false {
		t.Fatalf("Received code %s (expected 200) with body: %s and passcode %s", strconv.Itoa(statusCode), body, passcode)
	}
}

package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func makeReq(endpoint string, method string, json string, t *testing.T) (int, string) {
	url := fmt.Sprintf("http://backend:8080%s", endpoint)
	fmt.Printf("Recieved method %s for URL %s", method, url)

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

	if statusCode != 401 || strings.Compare(body, "invalid username") != 0 {
		t.Fatalf("Received code %s (expected 401) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidUsernameButBadPassword(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "bad_pass"
	}`, t)

	if statusCode != 401 || strings.Compare(body, "invalid password") != 0 {
		t.Fatalf("Received code %s (expected 401) with error: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidUsernameAndPassword(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "kenton",
		"password": "password123"
	}`, t)

	if statusCode != 200 || strings.Compare(body, "session token somewhere here") != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestLoginWithValidButOddlyCapitalizedUsernameAndPassword(t *testing.T) {
	statusCode, body := makeReq("/login", "POST", `{
		"username": "KenToN",
		"password": "password123"
	}`, t)

	if statusCode != 200 || strings.Compare(body, "session token somewhere here") != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

// TOTP Tests
func TestTOTPFailsWhenUsingInvalidType(t *testing.T) {
	statusCode, body := makeReq("/mfa/add", "POST", `{
		"type": "bad_type"
	}`, t) // Normally available: totp, fido

	if statusCode != 200 || strings.Compare(body, "invalid mfa type") != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPCreationSuccessful(t *testing.T) {
	statusCode, body := makeReq("/mfa/add", "POST", `{
		"type": "totp"
	}`, t)

	match, _ := regexp.MatchString(`{"totp": ".+"}`, body)

	if statusCode != 200 || match == false {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

func TestTOTPFailsWhenTryingToCreateSecondOTP(t *testing.T) {
	statusCode, body := makeReq("/mfa/add", "POST", `{
		"type": "totp"
	}`, t)

	if statusCode != 200 || strings.Compare(body, "account already has totp") != 0 {
		t.Fatalf("Received code %s (expected 200) with body: %s", strconv.Itoa(statusCode), body)
	}
}

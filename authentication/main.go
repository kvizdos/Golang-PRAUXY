package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
)

func contains(s []reflect.Value, str string) bool {
	for _, v := range s {
		if v.String() == str {
			return true
		}
	}

	return false
}

type HandlerFunc func(http.ResponseWriter, *http.Request, map[string]string)

type BodyField struct {
	name     string
	required bool
	//TODO: Regex
}

type RequestHandler struct {
	handler      HandlerFunc
	methodRoutes map[string][]BodyField
}

func (h *RequestHandler) handle(w http.ResponseWriter, r *http.Request) {
	if !contains(reflect.ValueOf(h.methodRoutes).MapKeys(), r.Method) {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var body map[string]string

	if len(h.methodRoutes[r.Method]) != 0 {
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}

		if err := json.Unmarshal([]byte(b), &body); err != nil {
			http.Error(w, "Bad request, missing body", http.StatusBadRequest)
			return
		}

		for _, field := range h.methodRoutes[r.Method] {
			_, exists := body[field.name]
			if field.required && !exists {
				http.Error(w, fmt.Sprintf("Missing body field %s", field.name), http.StatusBadRequest)
				return
			}
		}
	}

	w.Header().Add("Content-Type", "application/json")

	h.handler(w, r, body)
}

func main() {
	fmt.Println("Authentication backend started")

	registrationRoutes := make(map[string][]BodyField)
	registrationRoutes["POST"] = []BodyField{
		BodyField{
			name:     "username",
			required: true,
		},
		BodyField{
			name:     "email",
			required: false,
		},
		BodyField{
			name:     "password",
			required: true,
		},
	}
	RegistrationHandler := RequestHandler{
		handler:      RegisterUserHTTPWrapper,
		methodRoutes: registrationRoutes,
	}

	loginRoutes := make(map[string][]BodyField)
	loginRoutes["POST"] = []BodyField{
		BodyField{
			name:     "username",
			required: true,
		},
		BodyField{
			name:     "password",
			required: true,
		},
	}
	LoginHandler := RequestHandler{
		handler:      LoginUserHTTPWrapper,
		methodRoutes: loginRoutes,
	}

	http.HandleFunc("/register", RegistrationHandler.handle)
	http.HandleFunc("/login", LoginHandler.handle)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

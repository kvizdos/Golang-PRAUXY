package main

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"net/http"

	b64 "encoding/base64"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/pquerna/otp/totp"
)

type TotpInfo struct {
	secret string
	qr     string
}

type mfaInterface struct {
	TypeOfMFA string `bson:"type"` // Available: totp, fido
	Secret    string `bson:"secret"`
}

func getActiveMFA(username string) []mfaInterface {
	var result User

	conn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	return result.Multifactor
}

func createTOTP(username string) TotpInfo {
	totpOpts := totp.GenerateOpts{
		Issuer:      "PRAUXY",
		AccountName: username,
	}
	key, _ := totp.Generate(totpOpts)

	fmt.Println(username)

	var buf bytes.Buffer
	img, _ := key.Image(200, 200)
	png.Encode(&buf, img)

	totp := TotpInfo{
		secret: key.Secret(),
		qr:     b64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	dbTotp := mfaInterface{
		TypeOfMFA: "totp",
		Secret:    totp.secret,
	}

	filter := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"multifactor": dbTotp}}

	conn.UpdateOne(context.TODO(), filter, update)

	return totp
}

func AddMFAHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	print("Username: " + body["username"])

	ret := Response{}

	activeMfas := getActiveMFA(body["username"])

	totpEnabled := false
	// hardwareKeyEnabled := false
	// backuphwKeyEnabled := false

	for _, mfa := range activeMfas {
		if !totpEnabled {
			totpEnabled = mfa.TypeOfMFA == "totp"
		}
	}

	switch mfaType := body["type"]; mfaType {
	case "totp":
		if totpEnabled {
			ret.statusCode = 406
			ret.body = "totp already registered"
			break
		}
		totp := createTOTP(body["username"])
		ret.statusCode = 200
		ret.body = fmt.Sprintf(`{"qr": "%s"}`, totp.qr)
	default:
		ret.statusCode = 406
		ret.body = "invalid mfa type"
	}

	w.WriteHeader(ret.statusCode)
	fmt.Fprint(w, ret.body)

}

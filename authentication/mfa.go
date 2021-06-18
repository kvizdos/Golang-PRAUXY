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
	Secret string
	Qr     string
}

type mfaInterface struct {
	TypeOfMFA string `bson:"type"` // Available: totp, fido
	Secret    string `bson:"secret"`
}

type Mfas struct {
	totp     Mfa
	hardware Mfa
	backuphw Mfa
}

type Mfa struct {
	enabled bool
	secret  string
}

func getActiveMFA(username string) Mfas {
	var result User

	conn.FindOne(context.TODO(), User{
		Username: "kenton",
	}).Decode(&result)

	totp := Mfa{enabled: false}
	hardwareKey := Mfa{enabled: false}
	backupHardwareKey := Mfa{enabled: false}

	for _, mfa := range result.Multifactor {
		if !totp.enabled && mfa.TypeOfMFA == "totp" {
			totp.enabled = true
			totp.secret = mfa.Secret
		}
		if !hardwareKey.enabled {
			hardwareKey.enabled = mfa.TypeOfMFA == "hardware"
		}
		if !backupHardwareKey.enabled {
			backupHardwareKey.enabled = mfa.TypeOfMFA == "hardware-backup"
		}
	}

	return Mfas{
		totp:     totp,
		hardware: hardwareKey,
		backuphw: backupHardwareKey,
	}
}

func createTOTP(username string) TotpInfo {
	totpOpts := totp.GenerateOpts{
		Issuer:      "PRAUXY",
		AccountName: username,
	}
	key, _ := totp.Generate(totpOpts)

	var buf bytes.Buffer
	img, _ := key.Image(200, 200)
	png.Encode(&buf, img)

	totp := TotpInfo{
		Secret: key.Secret(),
		Qr:     b64.StdEncoding.EncodeToString(buf.Bytes()),
	}

	dbTotp := mfaInterface{
		TypeOfMFA: "totp",
		Secret:    totp.Secret,
	}

	filter := bson.M{"username": username}
	update := bson.M{"$push": bson.M{"multifactor": dbTotp}}

	conn.UpdateOne(context.TODO(), filter, update)

	return totp
}

func disableMFA(username string, mfaType string) {
	filter := bson.M{"username": username}
	update := bson.M{"$pull": bson.M{"multifactor": bson.M{"type": mfaType}}}

	conn.UpdateOne(context.TODO(), filter, update)
}

func verifyTotp(secret string, code string) bool {
	return totp.Validate(code, secret)
}

func MFAHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	switch method := r.Method; method {
	case "POST":
		AddMFAHTTPWrapper(w, r, body)
	case "DELETE":
		DeleteMFAHTTPWrapper(w, r, body)
	}
}

func MFAVerificationHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	ret := Response{}

	activeMfas := getActiveMFA(body["username"])

	switch mfaType := body["type"]; mfaType {
	case "totp":
		if !activeMfas.totp.enabled {
			ret.statusCode = 406
			ret.body = "totp not enabled"
			break
		}

		verified := verifyTotp(activeMfas.totp.secret, body["code"])

		if verified {
			ret.statusCode = 200
			ret.body = "session token here eventually"
		} else {
			ret.statusCode = 403
			ret.body = "invalid totp code"
		}
	default:
		ret.statusCode = 406
		ret.body = "invalid mfa type"
	}

	w.WriteHeader(ret.statusCode)
	fmt.Fprint(w, ret.body)
}

func AddMFAHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	ret := Response{}

	activeMfas := getActiveMFA(body["username"])

	switch mfaType := body["type"]; mfaType {
	case "totp":
		if activeMfas.totp.enabled {
			ret.statusCode = 406
			ret.body = "totp already registered"
			break
		}
		totp := createTOTP(body["username"])
		ret.statusCode = 200
		ret.body = fmt.Sprintf(`{"qr": "%s"}`, totp.Qr)
	default:
		ret.statusCode = 406
		ret.body = "invalid mfa type"
	}

	w.WriteHeader(ret.statusCode)
	fmt.Fprint(w, ret.body)
}

func DeleteMFAHTTPWrapper(w http.ResponseWriter, r *http.Request, body map[string]string) {
	ret := Response{}

	activeMfas := getActiveMFA(body["username"])

	switch mfaType := body["type"]; mfaType {
	case "totp":
		if !activeMfas.totp.enabled {
			ret.statusCode = 406
			ret.body = "totp not enabled"
			break
		}
		disableMFA(body["username"], "totp")

		ret.statusCode = 200
		ret.body = "totp disabled"
	default:
		ret.statusCode = 406
		ret.body = "invalid mfa type"
	}

	w.WriteHeader(ret.statusCode)
	fmt.Fprint(w, ret.body)
}

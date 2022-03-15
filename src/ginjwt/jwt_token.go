package ginjwt

import (
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

func GinJwtToken(userID string) (map[string]interface{}, bool) {

	mw := MwInitializer()

	// Create the token
	token := jwt.New(jwt.GetSigningMethod(mw.SigningAlgorithm))
	claims := token.Claims.(jwt.MapClaims)

	if mw.PayloadFunc != nil {
		for key, value := range mw.PayloadFunc(userID) {
			claims[key] = value
		}
	}

	expire := mw.TimeFunc().Add(mw.Timeout)
	claims["id"] = userID
	claims["exp"] = expire.Unix()
	claims["orig_iat"] = mw.TimeFunc().Unix()

	mapd := map[string]interface{}{"token": "", "expire": ""}

	tokenString, err := token.SignedString(mw.Key)
	if err != nil {
		return mapd, false
	}

	mapd = map[string]interface{}{"error": false, "token": tokenString, "expire": expire.Format(time.RFC3339)}

	return mapd, true
}

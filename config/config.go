package config

import (
	"os"
)

// JWTSignaturePrivateKey secret key
var JWTSignaturePrivateKey string = "key"

// CheckAdminUseridAndPass Check authenticity
func CheckAdminUseridAndPass(userid, password string) bool {
	uid := os.Getenv("XS_ADMIN_USERID")

	pass := os.Getenv("XS_ADMIN_PASSWORD")

	if uid == userid && pass == password {
		return true
	}
	return false
}

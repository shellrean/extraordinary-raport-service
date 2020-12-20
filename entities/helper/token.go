package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/shellrean/extraordinary-raport/domain"
)

func CreateToken(key string, user domain.User) (token string, err error) {
	// Create access token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = user.ID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	
	// Generate access token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	
	// Get token signed
	token, err = at.SignedString([]byte(key))
	if err != nil {
		return token, err
	}

	return
}
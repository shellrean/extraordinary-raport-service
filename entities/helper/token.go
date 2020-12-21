package helper

import (
	"github.com/dgrijalva/jwt-go"

	"github.com/shellrean/extraordinary-raport/domain"
)

func CreateAccessToken(key string, user domain.User, td *domain.TokenDetails) (err error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.ID
	atClaims["exp"] = td.AtExpires
	
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(key))
	if err != nil {
		return
	}
	return
}

func CreateRefreshToken(key string, user domain.User, td *domain.TokenDetails) (err error) {
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user.ID
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(key))
	if err != nil {
		return
	}
	return
}
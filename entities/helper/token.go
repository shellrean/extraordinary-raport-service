package helper

import (
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"

	"github.com/shellrean/extraordinary-raport/domain"
)

func GenerateTokenDetail(td *domain.Token) {
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
    td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
    td.AccessUuid = uuid.NewV4().String()
    td.RefreshUuid = uuid.NewV4().String()
}

func CreateAccessToken(key string, user domain.User, td *domain.Token) (err error) {
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = user.ID
	atClaims["role"] = user.Role
	atClaims["exp"] = td.AtExpires
	
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(key))
	if err != nil {
		return domain.ErrSessDecode
	}
	return
}

func CreateFileAccessToken(key string, path string) (token string, err error) {
	atClaims := jwt.MapClaims{}
	atClaims["path"] = path
	atClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token,  err = at.SignedString([]byte(key))
	if err != nil {
		return "", domain.ErrSessDecode
	}
	return
}

func CreateRefreshToken(key string, user domain.User, td *domain.Token) (err error) {
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = user.ID
	rtClaims["exp"] = td.RtExpires

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(key))
	if err != nil {
		return domain.ErrSessDecode
	}
	return
}

func ExtractToken(bearer string) (res string) {
	str := strings.Split(bearer, " ")
	if len(str) == 2 {
		res = str[1]
		return
	}
	return
}

func VerifyToken(key string, tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrSessDecode
		}
		return []byte(key), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired{
			return nil, domain.ErrSessExpired
		}

		return nil, domain.ErrSessVerifation
	}
	return token, nil
}

func TokenValid(token *jwt.Token) error {
	if _, ok := token.Claims.(jwt.Claims); !ok || !token.Valid {
		return domain.ErrSessExpired
	}
	return nil
}

func ExtractTokenMetadata(token *jwt.Token) map[string]interface{}{
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims
	}
	return nil
}
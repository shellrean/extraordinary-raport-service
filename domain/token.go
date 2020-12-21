package domain

type TokenDetails struct {
	AccessToken		string
	RefreshToken	string
	AccessUuid		string
	RefreshUuid		string
	AtExpires		int64
	RtExpires		int64
}

type DTOTokenResponse struct {
	AccessToken 	string 		`json:"access_token"`
	RefreshToken	string 		`json:"refresh_token"`
}
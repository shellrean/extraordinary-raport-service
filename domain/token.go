package domain

type Token struct {
	AccessToken		string
	RefreshToken	string
	AccessUuid		string
	RefreshUuid		string
	AtExpires		int64
	RtExpires		int64
}
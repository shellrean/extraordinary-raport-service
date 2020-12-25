package helper

import (
	"encoding/base64"
	"strconv"
)

func DecodeCursor(encodedString string) (int64, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return 0, err
	}
	curS := string(byt)
	var cursor int
	cursor, err = strconv.Atoi(curS)
	if err != nil {
		return 0, err
	}
	return int64(cursor), err
}

func EncodeCursor(cursor int64) string {
	cursorS := strconv.Itoa(int(cursor))

	return base64.StdEncoding.EncodeToString([]byte(cursorS))
}
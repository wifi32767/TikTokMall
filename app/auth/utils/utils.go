package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/dromara/dongle"
)

func GenerateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

type Token struct {
	Userid int32
	Salt   string
}

func GenerateToken(uid int32) (string, error) {
	salt, err := GenerateRandomString(16)
	if err != nil {
		return "", err
	}
	token := Token{
		Userid: uid,
		Salt:   salt,
	}
	jsonBytes, err := json.Marshal(token)
	if err != nil {
		return "", err
	}
	return dongle.Encode.FromBytes(jsonBytes).ByBase64().ToString(), nil
}

func ParseToken(token string) (*Token, error) {
	jsonBytes := dongle.Decode.FromString(token).ByBase64().ToBytes()
	var t Token
	err := json.Unmarshal(jsonBytes, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

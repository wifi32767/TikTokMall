package utils

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"

	"github.com/dromara/dongle"
)

// 生成随机字符串
func GenerateRandomString(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// 用Userid和随机的盐值生成一个token
// 盐值作为标记，防止伪造
// 同时可以轻易获取token对应的userid
type Token struct {
	Userid uint32
	Salt   string
}

func GenerateToken(uid uint32) (string, error) {
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

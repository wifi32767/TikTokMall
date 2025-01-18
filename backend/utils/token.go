package utils

import (
	"encoding/json"

	"github.com/dromara/dongle"
)

type Token struct {
	Userid int32
	Salt   string
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

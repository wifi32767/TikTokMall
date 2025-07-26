package message

import (
	"context"
	"encoding/json"
)

type AsyncMessage struct {
	Ctx         context.Context
	ServiceName string
	MethodName  string
	Body        any
}

func Encode(ctx context.Context, serviceName, methodName string, body any) ([]byte, error) {
	return json.Marshal(AsyncMessage{
		Ctx:         ctx,
		ServiceName: serviceName,
		MethodName:  methodName,
		Body:        body,
	})
}

func Decode(s []byte) (AsyncMessage, error) {
	var msg AsyncMessage
	err := json.Unmarshal(s, &msg)
	return msg, err
}

package message

import "context"

type AsyncMessage struct {
	Head string
	Body []byte
	Ctx  context.Context
}

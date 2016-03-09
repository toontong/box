package ping

import (
	"golang.org/x/net/context"
)

type PingService struct{}

func (_ PingService) Ping(_ context.Context, msg *MsgPing) (*Pong, error) {
	p := new(Pong)
	p.Msg = msg.Msg
	return p, nil
}

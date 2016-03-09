package gateway

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	pb "github.com/toontong/box/proto/nameserver"
	"github.com/toontong/box/proto/ping"
)

type Gateway struct {
	ping.PingService
}

func NewGateway() *Gateway {
	g := new(Gateway)
	return g
}

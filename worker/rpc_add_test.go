package worker

import (
	// "fmt"
	"math/rand"
	// "net"
	"testing"
	// "time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	pb "github.com/toontong/box/proto/worker"
)

const (
	Worker_Host = "127.0.0.1:8080"
)

func Test_RPC_ADD(t *testing.T) {
	var opts []grpc.DialOption
	// TODO:先不使用TLS
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(Worker_Host, opts...)
	if err != nil {
		t.Fatalf("Failed to dial[%s]: err=[%v].", Worker_Host, err)
	}

	defer conn.Close()
	client := pb.NewWrokerClient(conn)

	add := new(pb.AddReq)
	add.A = rand.Int63()
	add.B = rand.Int63()

	resp, err := client.Add(context.Background(), add)
	if err == nil && resp.Sum == add.A+add.B {
		log.Info("call rpc Success.")
	} else {
		t.Fatalf("Failed to call Worker[%v] RPC. err=[%v]", Worker_Host, err)
	}
}

package main

import (
	"flag"
	"fmt"
	"net"

	// "github.com/golang/protobuf/proto"
	// "golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	pb "github.com/toontong/box/proto/nameserver"

	"github.com/toontong/box/nameserver"
)

var (
	port = flag.Int("port", 20000, "The NameServer listen port")
)

func main() {
	log.Info("pageage import done, go into run main.")
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%d", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		panic(err)
	}

	var opts []grpc.ServerOption
	// if *tls {
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		grpclog.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterNameServiceServer(grpcServer, nameserver.NewNameServer())

	log.Infof("Serve in[%s]", addr)
	grpcServer.Serve(lis)
}

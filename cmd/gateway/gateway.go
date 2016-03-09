package main

import (
	"flag"
	"time"

	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	pb "github.com/toontong/box/proto/nameserver"

	"github.com/toontong/box/gateway"
)

var (
	nameserver = flag.String("nameserver", "127.0.0.1:20000", "NameServer ip:port,default[127.0.0.1:20000]")
	port       = flag.Int("port", 8080, "listen port,default[8080]")
	host       = flag.String("host", "0.0.0.0", "listen host,default[0.0.0.0]")
)

func main() {
	flag.Parse()

	wk := gateway.NewGateway(*host, *port)
	lis, err := wk.Listen()
	if err != nil {
		panic(err)
	}

	wk.JoinNameServer(*nameserver, 3*time.Second)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterWrokerServer(grpcServer, wk)
	log.Infof("Work going to Servr[%v:%v]", *host, *port)
	grpcServer.Serve(lis)
	wk.Stop()
}

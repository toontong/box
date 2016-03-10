package main

import (
	"flag"
	"time"

	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	"github.com/toontong/box/proto/ping"
	pb "github.com/toontong/box/proto/worker"
	"github.com/toontong/box/worker"
)

var (
	nameserver = flag.String("nameserver", "127.0.0.1:20000", "NameServer ip:port,default[127.0.0.1:20000]")
	port       = flag.Int("port", 40000, "listen port,default[40000]")
	host       = flag.String("host", "0.0.0.0", "listen host,default[0.0.0.0]")
)

func main() {
	flag.Parse()

	wk := worker.NewWoker()
	lis, err := wk.Listen(*host, *port)
	if err != nil {
		panic(err)
	}

	wk.JoinNameServer(*nameserver, 3*time.Second)
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)

	pb.RegisterWrokerServer(grpcServer, wk)
	ping.RegisterPingServiceServer(grpcServer, wk)

	log.Infof("Work going to Servr[%v:%v]", *host, *port)
	grpcServer.Serve(lis)
	wk.Stop()
}

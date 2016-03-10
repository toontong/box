package gateway

import (
	"fmt"
	"io"
	// "math/rand"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	"github.com/toontong/box/nameserver"
	pb "github.com/toontong/box/proto/nameserver"
	// "github.com/toontong/box/proto/ping"

	"github.com/toontong/box/libs/StreamServer"
)

type Gateway struct {
	// ping.PingService
	tcpSvr  *StreamServer.TCPServer
	nameCli pb.NameServiceClient

	nameServAddr string

	wkMapper map[uint64]*worker

	listenAddr string
	stop       bool
}

type worker struct {
	pb.Worker
	Alive bool
	RTT   time.Duration //RTT mean: round-trip time
}

func (g *Gateway) autoRefreshWorkerlist(period time.Duration) {
	go func() {
		var opts []grpc.DialOption
		// TODO:先不使用TLS
		opts = append(opts, grpc.WithInsecure())
		log.Infof("Will Refresh Workerlist from NameServer[%v] every[%v].", g.nameServAddr, period)

		for !g.stop {
			if g.nameCli == nil {
				conn, err := grpc.Dial(g.nameServAddr, opts...)
				if err != nil {
					log.Fatalf("Failed to dial[%s]: err=[%v],will retry after[%s]",
						g.nameServAddr, err, period)
					time.Sleep(period)
					continue
				}
				g.nameCli = pb.NewNameServiceClient(conn)
			}

			req := pb.Req{}

			stream, err := g.nameCli.ListWorkers(context.Background(), &req)
			if err != nil {
				log.Fatalf("Failed to JoinNameServer[%s]: err=[%v],will retry after[%s]",
					g.nameServAddr, err)
				time.Sleep(period)
				continue
			}
			for {
				wk, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					log.Warnf("%v.ListWorkers(_) err=[%s]", g.nameCli, err)
					continue
				}
				g.pushWorker(wk)
			}

			g.refreshWorkerRTT()
			time.Sleep(period)
		}
	}()
}

func (g *Gateway) pushWorker(wk *pb.Worker) {
	if w, ok := g.wkMapper[wk.WorkerId]; ok {
		w.Worker = *wk
	} else {
		w := worker{
			Worker: *wk,
		}
		g.wkMapper[wk.WorkerId] = &w
	}
}

func (g *Gateway) refreshWorkerRTT() {
	for _, wk := range g.wkMapper {
		rtt, err := nameserver.PingWorker(&wk.Worker)
		if err == nil {
			wk.RTT = rtt
			wk.Alive = true
		} else {
			wk.Alive = false
		}
	}
}

func NewGateway(nameServAddr string,
	listenHost string, listenPort int) *Gateway {

	g := new(Gateway)
	if nameServAddr == "" {
		panic("nameServAddr can not be empty.")
	}
	g.nameServAddr = nameServAddr

	g.listenAddr = fmt.Sprintf("%v:%v", listenHost, listenPort)
	g.wkMapper = make(map[uint64]*worker, 8)
	return g
}

func (g *Gateway) Serve() error {
	g.tcpSvr = StreamServer.NewTCPServer(g.listenAddr, g)
	g.autoRefreshWorkerlist(5 * time.Second)

	return g.tcpSvr.Serve()
}

func (g *Gateway) Stop() {
	g.stop = true
}

func (g Gateway) getAliveWorker() *worker {
	if len(g.wkMapper) == 0 {
		return nil
	}

	for _, w := range g.wkMapper {
		if w.Alive {
			return w
		}
	}
	return nil
}

// Do nothing on accept error.
func (g Gateway) OnAcceptError(err error) {
	log.Warnf("OnAcceptError:%v", err)
}

func (g *Gateway) OnConnection(conn net.Conn) {
	log.Info("Conn Create from[%s]", conn.RemoteAddr())
	wk := g.getAliveWorker()
	if wk == nil {
		log.Warnf("None Worker Alive.")
		conn.Write([]byte("None worker alive.\n"))
		conn.Close()
		return
	}
	var err error
	var wkConn net.Conn
	var nTry = 3
	for nTry > 0 {
		nTry--
		wkConn, err = net.Dial("tcp", wk.ListenAddr)
		if err != nil {
			log.Errorf("Failed to Dial(%v) err=[%s] nTry=[%d]", wk.ListenAddr, err, nTry)
			continue
		} else {
			break
		}
	}
	if wkConn == nil {
		log.Errorf("Failed Dial(%v) after retry.", wk.ListenAddr)
		return
	}
	go inbound(conn, wkConn)
	outbound(conn, wkConn)
}
func inbound(in net.Conn, out net.Conn) {
	_, err := io.Copy(in, out)
	if err != nil {
		in.Close()
		out.Close()
		log.Infof("inbound: Close(%s) and Close(%s)",
			in.RemoteAddr(), out.RemoteAddr())
	}
}

func outbound(in net.Conn, out net.Conn) {
	_, err := io.Copy(out, in)
	if err != nil {
		in.Close()
		out.Close()
		log.Infof("outbound: Close(%s) and Close(%s)",
			in.RemoteAddr(), out.RemoteAddr())
	}
}

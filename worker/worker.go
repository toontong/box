package worker

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	"github.com/toontong/box/proto/nameserver"
	"github.com/toontong/box/proto/ping"
	pb "github.com/toontong/box/proto/worker"
)

type IWorker interface {
	JoinNameServer(nameServerAddr string, period time.Duration) error
	Stop()
	Listen(host string, port int) error
}

type Worker struct {
	// 继承ping/pong RPC接口
	ping.PingService

	host      string
	port      int
	cpuUsage  float64
	currConn  uint32
	closeConn uint32

	workerId uint64

	stop  bool
	start time.Time
}

func NewWoker() *Worker {
	w := new(Worker)

	w.workerId = 0 // 由nameserver分配
	w.cpuUsage = 0.08
	w.stop = false
	w.start = time.Now()
	return w
}

func (w Worker) Add(_ context.Context, req *pb.AddReq) (*pb.AddResp, error) {
	resp := new(pb.AddResp)
	resp.Sum = req.A + req.B
	return resp, nil
}

func (w *Worker) Listen(host string, port int) (net.Listener, error) {
	w.host = host
	w.port = port
	addr := fmt.Sprintf("%s:%d", w.host, w.port)

	//TODO: just tcp4 ?
	lis, err := net.Listen("tcp4", addr)
	if err != nil {
		log.Fatalf("Failed to listen on[%v]: err=[%v]", addr, err)
		return nil, err
	}
	return lis, nil
}

func (w *Worker) Stop() { w.stop = true }

func (w *Worker) JoinNameServer(nameServerAddr string, period time.Duration) error {
	go func() {
		var opts []grpc.DialOption

		// TODO:先不使用TLS
		opts = append(opts, grpc.WithInsecure())
		// var conn grpc.ClientConn
		// var err error

		log.Infof("Will Report to NameServer[%v] every[%v].", nameServerAddr, period)

		for !w.stop {
			conn, err := grpc.Dial(nameServerAddr, opts...)
			if err != nil {
				log.Fatalf("Failed to dial[%s]: err=[%v],will retry after[%s]",
					nameServerAddr, err, period)
				time.Sleep(period)
				continue
			}

			client := nameserver.NewNameServiceClient(conn)

			req := nameserver.JoinReq{
				WorkerId:        w.workerId,
				Host:            w.host,
				Port:            int32(w.port),
				CurrConnection:  0,
				CloseConnection: w.getCloseConn(),
				CpuUsage:        0.1,
			}
			resp, err := client.WorkerJoin(context.Background(), &req)
			if err != nil {
				log.Fatalf("Failed to JoinNameServer[%s]: err=[%v],will retry after[%s]",
					nameServerAddr, err)
				continue
			}
			conn.Close()
			if resp.Success {
				log.Infof("Join NameSrv[%s] Success got WorkerId=[%v]",
					nameServerAddr, resp.WorkerId)
				w.workerId = resp.WorkerId
			} else {
				log.Infof("Failed to Join NameSrv errMsg=[%s] WorkerId=[%v],will retry after[%s]",
					resp.ErrMsg, resp.WorkerId, period)
			}

			time.Sleep(period)
		}
	}()
	if nameServerAddr == "" {
		return fmt.Errorf("NameServer can not be empty.")
	}
	if period < time.Second {
		return fmt.Errorf("period can not less then 1 second.")
	}
	return nil
}

func (w *Worker) getCloseConn() uint32 {
	//TODO this
	w.closeConn++
	return w.closeConn
}

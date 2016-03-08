package worker

import (
	"fmt"
	"net"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/toontong/box/libs/log"
	"github.com/toontong/box/proto/nameserver"
	pb "github.com/toontong/box/proto/worker"
)

type IWorker interface {
	Report2NameServer(nameServerAddr string, period time.Duration) error
	Stop()
	Listen() error
}

type Worker struct {
	host      string
	port      int
	cpuUsage  float64
	currConn  uint32
	closeConn uint32

	stop  bool
	start time.Time
}

func NewWoker(host string, port int) *Worker {
	w := new(Worker)
	w.host = host
	w.port = port
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

func (w *Worker) Listen() (net.Listener, error) {
	addr := fmt.Sprintf("%s:%d", w.host, w.port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on[%v]: err=[%v]", addr, err)
		return nil, err
	}
	return lis, nil
}

func (w *Worker) Stop() { w.stop = true }

func (w *Worker) Report2NameServer(nameServerAddr string, period time.Duration) error {
	go func() {
		var opts []grpc.DialOption

		// 先不使用TLS
		opts = append(opts, grpc.WithInsecure())

		var nClose uint32 = 0
		log.Infof("Will Report to NameServer[%v] every[%v].", nameServerAddr, period)
		for !w.stop {
			conn, err := grpc.Dial(nameServerAddr, opts...)
			if err != nil {
				log.Fatalf("fail to dial[%s]: err=[%v] ,will continue try after[%s]",
					nameServerAddr, err, period)
				time.Sleep(period)
				continue
			}
			defer conn.Close()
			client := nameserver.NewNameServiceClient(conn)

			req := nameserver.JoinReq{
				Host:            w.host,
				Port:            w.port,
				CurrConnection:  0,
				CloseConnection: nClose,
				CpuUsage:        0.1,
			}
			resp, err := client.WorkerJoin(context.Background(), &req)
			if err != nil {
				log.Errorf("Failed to Report2NameServer report-thread existed. err=[%s]", err)
				return
			}
			if resp.Success {
				log.Infof("Join NameSrv Success[%v] got WorkerId=[%v]",
					resp.Success, resp.WorkerId)
			} else {
				log.Infof("Failed to Join NameSrv errMsg=[%s] WorkerId=[%v]",
					resp.ErrMsg, resp.WorkerId)
			}
			nClose++
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

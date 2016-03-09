package nameserver

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	// "github.com/golang/protobuf/proto"
	"golang.org/x/net/context"

	"github.com/toontong/box/libs/log"

	pb "github.com/toontong/box/proto/nameserver"
	"github.com/toontong/box/proto/ping"
)

type nameServer struct {
	ping.PingService

	wkMaper map[uint64]*pb.Worker
}

func NewNameServer() *nameServer {
	s := new(nameServer)
	s.wkMaper = make(map[uint64]*pb.Worker, 8)
	return s
}

// 每5分钟向服务器(NameService)报活
func (s *nameServer) WorkerJoin(_ context.Context, req *pb.JoinReq) (*pb.JoinResp, error) {
	log.Infof("Recv Worker[%v:%v] join request.", req.Host, req.Port)

	var wk *pb.Worker
	if req.WorkerId != 0 {
		if w, ok := s.wkMaper[req.WorkerId]; ok {
			wk = w
			wk.WorkerId = req.WorkerId
		}
	}
	if wk == nil {
		wk = new(pb.Worker)
		wk.WorkerId = s.genWorkerId(req.Host, req.Port)
	}

	wk.ListenAddr = fmt.Sprintf("%v:%v", req.Host, req.Port)
	wk.CurrConnection = req.CurrConnection
	wk.CloseConnection = req.CloseConnection
	wk.CpuUsage = req.CpuUsage
	wk.Version = wk.Version
	wk.LastAlive = time.Now().UnixNano()

	s.join(wk)

	resp := new(pb.JoinResp)
	resp.Success = true
	resp.WorkerId = wk.WorkerId

	return resp, nil

}

//
func (s nameServer) ListWorkers(_ *pb.Req, stream pb.NameService_ListWorkersServer) (err error) {
	log.Infof("List-Workers.len=[%s]", len(s.wkMaper))
	for wId, wk := range s.wkMaper {
		if err = stream.Send(wk); err != nil {
			log.Warnf("failed to Send Worker(%s).id=[%s]", wk, wId)
			return err
		}
	}
	return nil
}

// Convert net.IP to int64
func inet_aton(ip string) int64 {
	bits := strings.Split(ip, ".")

	b0, _ := strconv.Atoi(bits[0])
	b1, _ := strconv.Atoi(bits[1])
	b2, _ := strconv.Atoi(bits[2])
	b3, _ := strconv.Atoi(bits[3])

	var sum int64

	sum += int64(b0) << 24
	sum += int64(b1) << 16
	sum += int64(b2) << 8
	sum += int64(b3)

	return sum
}

func (s nameServer) genWorkerId(host string, port int32) uint64 {
	sum := inet_aton(host)
	id := uint64(sum<<32 + int64(port))
	for {
		if _, ok := s.wkMaper[id]; !ok {
			return id
		}
		// 如果 id < 2^32 表示随机生成
		id = uint64(rand.Int31())
		log.Infof("Gen WorkerId by Rand()->host[%v:%v]->id[%v]",
			host, port, id)
	}
}

func (s *nameServer) join(wk *pb.Worker) {
	s.wkMaper[wk.WorkerId] = wk
	//TODO: tell  the gateway-server
}

package daemon

import (
	"context"
	"net"
	"sync"
	"runtime"

	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	test "github.com/import-yuefeng/BGPParser/pb/test"
	task "github.com/import-yuefeng/BGPParser/pb/task"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	PORT = ":2048"
)

type server struct{}

func Daemon() {
	log.Info("hello, now is daemon mode")

	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	test.RegisterGreeterServer(s, &server{})
	task.RegisterBGPTaskerServer(s, &server{})
	log.Info("start gRPC service...")
	s.Serve(lis)
}



func (s *server) SayHello(ctx context.Context, in *test.HelloRequest) (*test.HelloReply, error) {
	log.Infoln("request: ", in.Name)
	return &test.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) AddRawParse(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	log.Infoln("add raw-bgp parse task: ", in.Path)
	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) AddBGPParse(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	md := &MetaData{
		AsPathMap: make(map[string]*analysis.BGPInfo),
		rw: new(sync.RWMutex),
	}
	log.Infoln("add bgp parse task:", in.Path)

	md.parseBGPData(in.Path, runtime.NumCPU())

	return &task.TaskReply{Message: "Success"}, nil
}


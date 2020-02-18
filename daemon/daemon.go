package daemon

import (
	"context"
	"io"
	"net"
	"os"
	"runtime"
	"sync"

	task "github.com/import-yuefeng/BGPParser/pb/task"
	test "github.com/import-yuefeng/BGPParser/pb/test"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	PORT = ":2048"
)

type server struct{}

var (
	root *analysis.BGPBST
	md   *MetaData
)

func Daemon() {
	log.Info("hello, now is daemon mode")

	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// if *logPath != "" {
	lf, err := os.OpenFile("./analysis.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
	if err != nil {
		log.Errorf("Unable to open log file for writing: %s", err)
	} else {
		log.SetOutput(io.MultiWriter(lf, os.Stdout))
	}
	// }

	// go func() {
	// 	log.Println(http.ListenAndServe("bgp-analyze.automesh.org:8000", nil))
	// }()

	s := grpc.NewServer()
	test.RegisterGreeterServer(s, &server{})
	task.RegisterBGPTaskerServer(s, &server{})
	task.RegisterAPIServer(s, &server{})
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
	md = &MetaData{
		AsPathMap: sync.Map{},
	}
	log.Infoln("add bgp parse task:", in.Path)

	root = md.parseBGPData(in.Path, runtime.NumCPU())

	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) SearchIP(ctx context.Context, in *task.IPAddr) (*task.SearchReply, error) {
	if root != nil {
		hashcode, err := root.Search(in.Ip)
		if err != nil {
			log.Infoln(err)
			return &task.SearchReply{Result: err.Error()}, nil
		}
		if t, ok := md.AsPathMap.Load(hashcode); ok {
			if res, ok := t.(*analysis.BGPInfo); ok {
				return &task.SearchReply{Result: res.Prefix[0]}, nil
			}
		}
	}
	return &task.SearchReply{Result: "Faild"}, nil

}

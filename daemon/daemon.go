package daemon

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"runtime"

	task "github.com/import-yuefeng/BGPParser/pb/task"
	test "github.com/import-yuefeng/BGPParser/pb/test"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	root     *analysis.BGPBST
	oldroot  *analysis.BGPBST
	md       *MetaData
	oldmd    *MetaData
	inUpdate bool
)

const (
	PORT = ":2048"
)

func (d *Daemon) Run() {
	log.Info("hello, now is daemon mode")

	lis, err := net.Listen("tcp", PORT)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if d.logPath != "" {
		lf, err := os.OpenFile("./analysis.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0640)
		if err != nil {
			log.Errorf("Unable to open log file for writing: %s", err)
		} else {
			log.SetOutput(io.MultiWriter(lf, os.Stdout))
		}
	} else {
		log.SetOutput(io.MultiWriter(os.Stdout))
	}

	go func() {
		log.Println(http.ListenAndServe("bgp-analyze.automesh.org:8000", nil))
	}()

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
	go func() {
		if inUpdate {
			return
		}
		inUpdate = true
		oldmd = md
		oldroot = root
		if md == nil {
			md = NewMetaData()
		} else {
			md.TaskList = make([][]*analysis.SimpleBGPInfo, 16)
		}
		log.Infoln("add bgp parse task:", in.Path)
		root = md.parseBGPData(in.Path, runtime.NumCPU())
		inUpdate = false
		log.Infoln("parse task end...")
		return
	}()
	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) SearchIP(ctx context.Context, in *task.IPAddr) (*task.SearchReply, error) {
	log.Infoln("search...", in.Ip)
	return s.searchByIP(in.Ip)
}

func (s *server) searchByIP(ip string) (*task.SearchReply, error) {
	if inUpdate {
		root = oldroot
		md = oldmd
	}
	if root != nil {
		hashcodeList, err := root.Search(ip)
		log.Warnln(hashcodeList)
		if err != nil || len(hashcodeList) == 0 {
			log.Infoln(err)
			return &task.SearchReply{Result: err.Error()}, nil
		}
		hashcode := hashcodeList[len(hashcodeList)-1]
		for _, v := range hashcodeList {
			if t, ok := md.AsPathMap.Load(v); ok {
				if res, ok := t.(*analysis.SimpleBGPInfo); ok {
					log.Infoln(res.Prefix)
				}
			}
		}
		if t, ok := md.AsPathMap.Load(hashcode); ok {
			if res, ok := t.(*analysis.SimpleBGPInfo); ok {
				return &task.SearchReply{Result: res.Prefix[0]}, nil
			}
		}
		return &task.SearchReply{Result: "Faild, not found."}, nil
	}
	return &task.SearchReply{Result: "Building iptree..."}, nil
}

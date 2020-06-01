// MIT License

// Copyright (c) 2019 Yuefeng Zhu

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package daemon

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"runtime"
	"strings"

	"github.com/import-yuefeng/BGPParser/utils"

	task "github.com/import-yuefeng/BGPParser/pb/task"
	test "github.com/import-yuefeng/BGPParser/pb/test"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	compress "github.com/import-yuefeng/BGPParser/tools/compress"
	marshal "github.com/import-yuefeng/BGPParser/tools/marshal"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	root *analysis.BGPBST
	md   *MetaData
)

const (
	PORT = ":2048"
)

func (d *Daemon) Run() *grpc.Server {
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
	// go func() {
	// 	log.Println(http.ListenAndServe("bgp-analyze.automesh.org:8000", nil))
	// }()
	s := grpc.NewServer()
	test.RegisterGreeterServer(s, &server{})
	task.RegisterBGPTaskerServer(s, &server{})
	task.RegisterAPIServer(s, &server{})
	log.Info("start gRPC service...")
	s.Serve(lis)
	return s
}

func (s *server) SayHello(ctx context.Context, in *test.HelloRequest) (*test.HelloReply, error) {
	log.Infoln("request: ", in.Name)
	return &test.HelloReply{Message: "Hello " + in.Name}, nil
}

func (s *server) AddRawParse(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	go addRawParse(in.Path)
	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) LoadIPTree(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	log.Infoln("load iptree file: ", in.Path)
	return loadIPTree(in.Path)
}

func (s *server) SaveIPTree(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	go saveIPTree(in.Path[0])
	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) AddBGPParse(ctx context.Context, in *task.FilePath) (*task.TaskReply, error) {
	go addBGPParse(in.Path)
	return &task.TaskReply{Message: "Success"}, nil
}

func (s *server) SearchIP(ctx context.Context, in *task.IPAddr) (*task.SearchReply, error) {
	log.Infoln("search: ", in.Ip)
	if !utils.IsIP(in.Ip) {
		return &task.SearchReply{Result: "Invalid IP"}, nil
	}
	return searchByIP(in.Ip)
}

func searchByIP(ip string) (*task.SearchReply, error) {
	if root != nil {
		prefixList, err := root.Search(ip)
		log.Warnln(prefixList)
		if err != nil && len(prefixList) == 0 {
			log.Infoln(err)
			return &task.SearchReply{Result: err.Error()}, nil
		}
		return &task.SearchReply{Result: prefixList[len(prefixList)-1]}, nil
	}
	return &task.SearchReply{Result: "Building iptree..."}, nil
}

func addBGPParse(paths []string) {
	if md == nil {
		md = NewMetaData()
	} else {
		md.TaskList = make([][]*analysis.BGPInfo, 16)
	}
	log.Infoln("add bgp parse task:", paths)
	root = md.ParseBGPData(paths, runtime.NumCPU())
	return
}

func loadIPTree(paths []string) (*task.TaskReply, error) {
	files := compress.Uncompress(paths[0], "tmpDir")
	for _, file := range files {
		root = marshal.Unmarshal(file)
		if err := os.Remove(file); err != nil {
			log.Fatalln(err)
			return &task.TaskReply{Message: "Fail"}, nil
		}
	}
	return &task.TaskReply{Message: "Success"}, nil
}

func addRawParse(paths []string) {
	log.Infoln("add raw-bgp parse task: ", paths)
	for _, file := range strings.Split(paths[0], " ") {
		file := file
		go parseRIBData(file)
	}
}

func saveIPTree(path string) {
	log.Infoln("start encoding iptree")
	root.EncodeIPTree()
	log.Infoln("save iptree to: ", path)
	marshal.Marshal(root, path)
	compress.Compress(path, fmt.Sprintf("%s.zip", path))
}

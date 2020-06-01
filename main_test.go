package main

import (
	"os"
	"testing"
	"time"

	client "github.com/import-yuefeng/BGPParser/cmd/client"
	server "github.com/import-yuefeng/BGPParser/daemon"
	analysis "github.com/import-yuefeng/BGPParser/tools/analysis"
	"google.golang.org/grpc"
)

var (
	root *analysis.BGPBST
	cli  *client.Client
)

func TestMain(m *testing.M) {
	server := setup()
	cli = client.NewClient("127.0.0.1:2048")
	code := m.Run()
	shutdown(server)
	os.Exit(code)
}

func setup() *grpc.Server {
	daemonCli := &server.Daemon{}
	s := daemonCli.Run()
	return s
}

func shutdown(server *grpc.Server) {
	root = nil
	server.GracefulStop()
	time.Sleep(10)
	return
}

func TestParseRIBData(t *testing.T) {
	defaultTestFile := "rib.20200409.0600.bz2"
	testData := []string{defaultTestFile}
	err := cli.AddRawParse(testData)
	if err != nil {
		t.Errorf("TestParseRIBData: %v", err)
	}
}

func TestParseBGPData(t *testing.T) {
	defaultTestFile := "rib.20200409.0600.bz2.txt"
	testData := []string{defaultTestFile}
	err := cli.AddBGPParse(testData)
	if err != nil {
		t.Errorf("TestParseBGPData: %v", err)
	}
}

func TestSearchIPTreeByIP(t *testing.T) {
	var testData map[string]string = map[string]string{
		"1.1.1.1":         "1.1.1.0/24",
		"8.8.8.8":         "8.8.8.0/24",
		"1.2.3.4":         "1.2.128.0/17",
		"114.114.114.114": "114.114.112.0/21",
		"4.3.2.1":         "4.0.0.0/9",
		"2.3.3.3":         "2.2.0.0/16",
	}
	for k, v := range testData {
		if err, res := cli.Search(ipaddr); err != nil || v != res {
			t.Errorf("TestSearchIPTreeByIP: %v", err)
		}
	}
}

func TestSaveIPTree(t *testing.T) {
	defaultSaveFile := "iptree"
	testData := []string{defaultSaveFile}
	err := cli.SaveIPTree(testData)
	if err != nil {
		t.Errorf("TestSaveIPTree: %v", err)
	}
}

func TestLoadIPTree(t *testing.T) {
	defaultSaveFile := "iptree.zip"
	testData := []string{defaultSaveFile}
	err := cli.LoadIPTree(testData)
	if err != nil {
		t.Errorf("TestLoadIPTree: %v", err)
	}
}

func BenchmarkA(b *testing.B) {

}
func BenchmarkB(b *testing.B) {

}

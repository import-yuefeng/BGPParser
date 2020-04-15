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

package client

import (
	task "github.com/import-yuefeng/BGPParser/pb/task"
	test "github.com/import-yuefeng/BGPParser/pb/test"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	conn    *grpc.ClientConn
	address string
}

func NewClient(address string) *Client {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &Client{
		conn:    conn,
		address: address,
	}
}

func (c *Client) ReConnect() {
	c.conn.Close()
	conn, err := grpc.Dial(c.address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	c.conn = conn
	return
}

func (c *Client) SayHello() {
	sayHello(c.conn)
}

func (c *Client) AddRawParse(filepath []string) {
	addRawParse(c.conn, filepath)
}

func (c *Client) AddBGPParse(filepath []string) {
	addBGPParse(c.conn, filepath)
}

func (c *Client) Search(ipaddr string) {
	search(c.conn, ipaddr)
}

func (c *Client) SaveIPTree(filepath []string) {
	saveIPTree(c.conn, filepath)
}

func (c *Client) LoadIPTree(filepath []string) {
	loadIPTree(c.conn, filepath)
}

func saveIPTree(conn *grpc.ClientConn, filepath []string) {
	c := task.NewBGPTaskerClient(conn)
	r, err := c.SaveIPTree(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Message)
}

func loadIPTree(conn *grpc.ClientConn, filepath []string) {
	c := task.NewBGPTaskerClient(conn)
	r, err := c.LoadIPTree(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Message)
}

func sayHello(conn *grpc.ClientConn) {
	c := test.NewGreeterClient(conn)
	name := "running..."
	r, err := c.SayHello(context.Background(), &test.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Message)
}

func addRawParse(conn *grpc.ClientConn, filepath []string) {
	c := task.NewBGPTaskerClient(conn)

	r, err := c.AddRawParse(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Message)
}

func addBGPParse(conn *grpc.ClientConn, filepath []string) {
	c := task.NewBGPTaskerClient(conn)

	r, err := c.AddBGPParse(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Message)
}

func search(conn *grpc.ClientConn, ip string) {
	c := task.NewAPIClient(conn)

	r, err := c.SearchIP(context.Background(), &task.IPAddr{Ip: ip})
	if err != nil {
		log.Fatalf("process error: ", err)
	}
	log.Println(r.Result)
}

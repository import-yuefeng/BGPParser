package client

import (
	test "github.com/import-yuefeng/BGPParser/pb/test"
	task "github.com/import-yuefeng/BGPParser/pb/task"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
	address string
}

func NewClient(address string) *Client {
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return &Client{
		conn: conn,
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

func (c *Client) AddRawParse(filepath string) {
	addRawParse(c.conn, filepath)
}

func (c *Client) AddBGPParse(filepath string) {
	addBGPParse(c.conn, filepath)
}

func sayHello(conn *grpc.ClientConn) {
	c := test.NewGreeterClient(conn)
	name := "running..."
	r, err := c.SayHello(context.Background(), &test.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.Message)
}

func addRawParse(conn *grpc.ClientConn, filepath string) {
	c := task.NewBGPTaskerClient(conn)

	r, err := c.AddRawParse(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.Message)
}

func addBGPParse(conn *grpc.ClientConn, filepath string) {
	c := task.NewBGPTaskerClient(conn)

	r, err := c.AddBGPParse(context.Background(), &task.FilePath{Path: filepath})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Println(r.Message)
}



package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "github.com/sdfwds4/test_go-zero_qps/proto"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/zrpc"
)

const name = "kevin"

var configFile = flag.String("f", "../etc/config.json", "the config file")

func main() {
	flag.Parse()

	fmt.Println("client start ...")

	start := time.Now()

	var rsp *pb.Response
	var err error

	var c zrpc.RpcClientConf
	conf.MustLoad(*configFile, &c)
	client := zrpc.MustNewClient(c)
	conn := client.Conn()
	for {
		// Make request
		greet := pb.NewGreeterClient(conn)
		rsp, err = greet.Greet(context.Background(), &pb.Request{
			Name: "kevin",
		})
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}

	fmt.Println("duration:", time.Since(start))

	fmt.Println("rsp.Greeting:", rsp.Greet)
}

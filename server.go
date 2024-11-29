package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/sdfwds4/test_go-zero_qps/proto"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/zrpc"

	"google.golang.org/grpc"
)

var counter int64

var configFile = flag.String("f", "../etc/config.json", "the config file")

type GreetServer struct {
	lock     sync.Mutex
	alive    bool
	downTime time.Time
}

func NewGreetServer() *GreetServer {
	return &GreetServer{
		alive: true,
	}
}

func (gs *GreetServer) Greet(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	// fmt.Println("=>", req)

	atomic.AddInt64(&counter, 1)

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	return &pb.Response{
		Greet: "hello from " + hostname,
	}, nil
}

func main() {
	flag.Parse()

	var c zrpc.RpcServerConf
	conf.MustLoad(*configFile, &c)
	logx.SetLevel(logx.ErrorLevel)
	logx.Disable()
	logx.Close()

	server := zrpc.MustNewServer(c, func(grpcServer *grpc.Server) {
		pb.RegisterGreeterServer(grpcServer, NewGreetServer())
	})
	// interceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// 	// st := time.Now()
	// 	resp, err = handler(ctx, req)
	// 	// log.Printf("method: %s time: %v\n", info.FullMethod, time.Since(st))
	// 	return resp, err
	// }

	// server.AddUnaryInterceptors(interceptor)

	go func() {
		var t = time.Now().UnixNano() / 1e6
		for {
			select {
			case <-time.After(time.Second * 5):
				now := time.Now().UnixNano() / 1e6
				v := atomic.SwapInt64(&counter, 0)
				log.Print("count: ", float64(v)/float64((now-t)/1000), "/s")
				t = now
			}
		}
	}()

	fmt.Println("server start ...")

	// Run server
	server.Start()
}

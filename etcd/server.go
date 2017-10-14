package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go-grpc-naming-demo/libs"
	"go-grpc-naming-demo/protos"

	"github.com/coreos/etcd/proxy/grpcproxy"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var HelloServerServiceName = "/etcd3_naming/grpc.hello"
var registTTL = 2

type HelloServer struct{}

func (h *HelloServer) SayHello(ctx context.Context, req *protos.HelloRequest) (*protos.HelloResponse, error) {
	response := &protos.HelloResponse{
		Reply: fmt.Sprintf("hello, %s", req.Greeting),
	}
	return response, nil
}

func main() {
	// parse flag
	port := flag.Int("port", 0, "the grpc server port")
	ipaddr := flag.String("ipaddr", "", "the grpc server ip addr")
	flag.Parse()

	// bind port
	listen, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		panic(err)
	}

	// start grpc server
	grpcServer := grpc.NewServer()
	protos.RegisterHelloServiceServer(grpcServer, &HelloServer{})
	reflection.Register(grpcServer)

	go func() {
		if err := grpcServer.Serve(listen); err != nil {
			panic(err)
		} else {
			fmt.Printf("grpc service start at port (%v)", *port)
		}
	}()

	// regist service in etcd
	addr := fmt.Sprintf("%s:%d", *ipaddr, *port)
	cli := libs.GetEtcdCli()
	grpcproxy.Register(cli, HelloServerServiceName, addr, registTTL)

	// wait for stop signal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	sign := <-signalChan
	fmt.Printf("receive signal (%v) ,grpc server will stop", sign)

	// graceful stop
	cli.Close()
	grpcServer.GracefulStop()
}

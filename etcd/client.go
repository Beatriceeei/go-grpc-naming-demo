package main

import (
	"fmt"
	"go-grpc-naming-demo/libs"
	"go-grpc-naming-demo/protos"
	"time"

	etcdnaming "github.com/coreos/etcd/clientv3/naming"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var HelloServerServiceName = "/etcd3_naming/grpc.hello"

type GrpcConnection interface {
	getGrpcConnection() *grpc.ClientConn
	connect(opt ...ConnOptions)
}
type ConnOptions func(conn GrpcConnection)

type EtcdResolverConnection struct {
	grpcConn    *grpc.ClientConn
	serviceName string
}

func (e *EtcdResolverConnection) connect(opts ...ConnOptions) {
	for _, opt := range opts {
		opt(e)
	}

	r := &etcdnaming.GRPCResolver{Client: libs.GetEtcdCli()}
	b := grpc.RoundRobin(r)
	ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	defer cancel()
	conn, gerr := grpc.DialContext(
		ctx,
		e.serviceName,
		grpc.WithInsecure(),
		grpc.WithBalancer(b),
		grpc.WithTimeout(time.Second*5),
		grpc.WithBlock(),
	)
	if gerr != nil {
		fmt.Printf("dial service(%s) by etcd resolver server error (%v)", e.serviceName, gerr.Error())
		panic(gerr)
	}
	e.grpcConn = conn

}

func (e *EtcdResolverConnection) getGrpcConnection() *grpc.ClientConn {
	return e.grpcConn
}

func main() {
	conn := &EtcdResolverConnection{
		serviceName: HelloServerServiceName,
	}
	conn.connect()

	for i := 0; i < 100; i++ {
		request := &protos.HelloRequest{Greeting: fmt.Sprintf("%d", i)}

		client := protos.NewHelloServiceClient(conn.getGrpcConnection())

		resp, err := client.SayHello(context.Background(), request)
		fmt.Printf("resp: %+v, err: %+v\n", resp, err)
		time.Sleep(time.Second)
	}

}

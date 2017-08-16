package connclose

import (
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/testutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func TestMain(m *testing.M) {
	m.Run()
	testutil.CheckLeakedGoroutine()
}

func TestGreeter(t *testing.T) {
	defer testutil.AfterTest(t)

	eps := make([]string, 3)
	for i := 0; i < 3; i++ {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		addr := ln.Addr().String()
		eps[i] = addr

		srv := grpc.NewServer()

		fmt.Println("serving:", addr)
		go srv.Serve(ln)

		s := &server{g: srv, l: ln}
		RegisterGreeterServer(srv, s)
	}

	test(t, eps)
}

type server struct {
	g *grpc.Server
	l net.Listener
}

func (s *server) Greet(stream Greeter_GreetServer) error {
	_, err := stream.Recv()
	return err
}

// stream, err := client.Hello(context.Background())
// if err != nil {
// 	panic(err)
// }

func test(t *testing.T, eps []string) {
	conn, err := grpc.Dial(
		eps[0],
		grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    time.Second,
			Timeout: time.Second,
		}),
		grpc.WithBalancer(clientv3.NewSimpleBalancer(eps)),
	)
	if err != nil {
		panic(err)
	}
	client := NewGreeterClient(conn)
	_ = client

	fmt.Println("waiting...")
	select {}
}

/*
add this to etcd/clientv3/balancer.go

func NewSimpleBalancer(eps []string) *simpleBalancer {
	return newSimpleBalancer(eps)
}


=== RUN   TestGreeter
serving: 127.0.0.1:60473
serving: 127.0.0.1:60474
serving: 127.0.0.1:60475
17:54:58.850108 [C] Notify!
17:54:58.850286 [C] Notify!
17:54:58.850820 [C] Up! 127.0.0.1:60473
waiting...
17:54:58.851080 [D] resetTransport(false): 127.0.0.1:60474 grpc: the connection is closing *errors.errorString
17:54:58.851184 [D] resetTransport(false): 127.0.0.1:60475 grpc: the connection is closing *errors.errorString
*/

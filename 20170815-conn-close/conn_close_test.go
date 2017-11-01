package connclose

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/pkg/testutil"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

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
		t.Fatal(err)
	}
	handleInterrupts(conn)

	fmt.Println("waiting after dial...")
	select {}
}

type server struct {
	g *grpc.Server
	l net.Listener
}

func (s *server) Greet(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	return nil, nil
}

func handleInterrupts(conn *grpc.ClientConn) {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-notifier
		fmt.Printf("received %q signal, shutting down...\n", sig)

		fmt.Printf("closing %q\n", conn.GetState())
		err := conn.Close()
		fmt.Println("conn.Close error:", err)

		signal.Stop(notifier)

		pid := syscall.Getpid()
		// exit directly if it is the "init" process, since the kernel will not help to kill pid 1.
		if pid == 1 {
			os.Exit(0)
		}
		syscall.Kill(pid, sig.(syscall.Signal))
	}()
}

func TestMain(m *testing.M) {
	m.Run()
	testutil.CheckLeakedGoroutine()
}

/*
add this to etcd/clientv3/balancer.go

func NewSimpleBalancer(eps []string) *simpleBalancer {
	return newSimpleBalancer(eps)
}

=== RUN   TestGreeter
serving: 127.0.0.1:61019
serving: 127.0.0.1:61020
serving: 127.0.0.1:61021
18:13:05.054162 [C] Notify!
18:13:05.054351 [C] Notify!
18:13:05.054919 [C] Up! 127.0.0.1:61019
18:13:05.055334 [D] resetTransport(false): 127.0.0.1:61020 grpc: the connection is closing *errors.errorString
18:13:05.055610 [D] resetTransport(false): 127.0.0.1:61021 grpc: the connection is closing *errors.errorString
waiting after dial...
^Creceived "interrupt" signal, shutting down...
closing "READY"
conn.Close error: <nil>
*/

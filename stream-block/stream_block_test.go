package streamblock

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/coreos/etcd/pkg/testutil"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	m.Run()
	testutil.CheckLeakedGoroutine()
}
func TestStreamBlock(t *testing.T) {
	defer testutil.AfterTest(t)

	// start server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	fmt.Println("serving:", addr)

	srv := grpc.NewServer()

	reqN := 10
	ch := make(chan string, reqN)
	RegisterElectionServer(srv, &testElectionServer{ch})

	go func() {
		srv.Serve(ln)
	}()

	cli, err := newClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(reqN)
	go func() {
		for i := 0; i < reqN; i++ {
			go func() {
				defer wg.Done()

				fmt.Println("Recv 1")
				_, err = cli.stream.Recv()
				fmt.Println("Recv 2")
				if err != nil {
					if strings.Contains(err.Error(), "transport is closing") {
						return
					}
					t.Fatal(err)
				}
			}()
		}
	}()

	go func() {
		time.Sleep(2 * time.Millisecond)
		for i := 0; i < reqN; i++ {
			ch <- "hello"
		}
	}()

	time.Sleep(2 * time.Millisecond)

	fmt.Println("Closing client!")
	cli.cancel()
	cli.conn.Close()
	fmt.Println("Closed client!")

	wg.Wait()

	srv.GracefulStop()
}

type client struct {
	conn   *grpc.ClientConn
	stream Election_ObserveClient
	ctx    context.Context
	cancel func()
}

func newClient(addr string) (*client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	stream, err := NewElectionClient(conn).Observe(ctx, &LeaderRequest{})
	if err != nil {
		conn.Close()
		cancel()
		return nil, err
	}
	return &client{conn, stream, ctx, cancel}, nil
}

type testElectionServer struct {
	ch chan string
}

func (srv *testElectionServer) Observe(_ *LeaderRequest, stream Election_ObserveServer) error {
	for {
		select {
		case v := <-srv.ch:
			if err := stream.Send(&ObserveResponse{Data: v + " (ack)"}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			return stream.Context().Err()
		}
	}
}

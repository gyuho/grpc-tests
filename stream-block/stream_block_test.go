package streamblock

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"
	"time"

	"strings"

	"github.com/coreos/etcd/pkg/testutil"
	"google.golang.org/grpc"
)

func TestMain(m *testing.M) {
	v := m.Run()
	if v == 0 && testutil.CheckLeakedGoroutine() {
		os.Exit(1)
	}
	os.Exit(v)
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

	ch := make(chan string, 1)
	RegisterElectionServer(srv, &testElectionServer{ch})

	go func() {
		srv.Serve(ln)
	}()
	defer srv.GracefulStop()

	cli, err := newClient(addr)
	if err != nil {
		t.Fatal(err)
	}

	stopc, donec := make(chan struct{}), make(chan struct{})
	go func() {
		for {
			select {
			case ch <- "hello":
			case <-stopc:
				close(donec)
				return
			}
		}
	}()

	go func() {
		for {
			_, err = cli.Recv()
			if err != nil {
				if strings.Contains(err.Error(), "transport is closing") {
					return
				}
				t.Fatal(err)
			}
		}
	}()

	time.Sleep(2 * time.Second)
	cli.Close()
	close(stopc)
	<-donec
}

type client struct {
	conn   *grpc.ClientConn
	stream Election_ObserveClient
}

func newClient(addr string) (*client, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	stream, err := NewElectionClient(conn).Observe(context.Background(), &LeaderRequest{})
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &client{conn, stream}, nil
}

func (c *client) Close() error { return c.conn.Close() }

func (c *client) Recv() (string, error) {
	rsp, err := c.stream.Recv()
	if err != nil {
		return "", err
	}
	return rsp.Data, nil
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
		time.Sleep(time.Second)
	}
}

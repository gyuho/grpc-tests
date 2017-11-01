package streamblock

import (
	"context"
	"fmt"
	"net"
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
func TestElectionBlock(t *testing.T) {
	defer testutil.AfterTest(t)

	// start server
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	fmt.Println("serving:", addr)

	srv := grpc.NewServer()

	go srv.Serve(ln)
	s := &server{g: srv, l: ln}
	RegisterElectionServer(srv, s)

	for {
		test(t, addr)
	}
}

type server struct {
	g *grpc.Server
	l net.Listener
}

func (s *server) Observe(stream Election_ObserveServer) error {
	_, err := stream.Recv()
	return err
}

func test(t *testing.T, endpoint string) {
	conn, err := grpc.Dial(endpoint, grpc.WithTimeout(time.Millisecond*100), grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := NewElectionClient(conn)
	stream, err := client.Observe(context.Background())
	if err != nil {
		panic(err)
	}
	c := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 500)
		fmt.Println(time.Now(), "Calling Close")
		conn.Close()
		t1 := time.Now()
		select {
		case <-time.After(time.Second):
			panic("Didn't close the connection within a timely manner")
		case <-c:
		}
		fmt.Println("Close took", time.Since(t1))

	}()
	fmt.Println(time.Now(), "Calling Recv")
	_, err = stream.Recv()
	fmt.Println(err)
	c <- struct{}{}
	wg.Wait()
}

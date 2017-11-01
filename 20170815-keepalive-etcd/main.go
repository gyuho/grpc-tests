package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/compactor"
	"github.com/coreos/etcd/embed"

	"github.com/golang/glog"
)

/*
go run main.go -logtostderr
*/

func init() {
	flag.Parse()
}

func main() {
	dataDir := filepath.Join(os.TempDir(), "etcd-test-data-dir")
	os.RemoveAll(dataDir)
	defer os.RemoveAll(dataDir)

	handleInterrupts(dataDir)
	start(dataDir)

	glog.Info("waiting... 1")
	time.Sleep(5 * time.Second)

	ccfg := clientv3.Config{
		Endpoints:            []string{"localhost:2379", "localhost:2381"},
		DialKeepAliveTime:    time.Second,
		DialKeepAliveTimeout: time.Second,
	}
	cli, err := clientv3.New(ccfg)
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	glog.Info("waiting... 2")
	select {}
}

func start(dataDir string) {
	cfgs := make([]*embed.Config, 3)
	iss := make([]string, 3)
	for i := range cfgs {
		cport := 2379 + 2*i
		pport := 2379 + 2*i + 1

		cfg := embed.NewConfig()
		cfg.ClusterState = embed.ClusterStateFlagNew

		cfg.Name = fmt.Sprintf("etcd-queue-%d", i)
		cfg.Dir = filepath.Join(dataDir, cfg.Name)

		curl := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", cport)}
		cfg.ACUrls, cfg.LCUrls = []url.URL{curl}, []url.URL{curl}

		purl := url.URL{Scheme: "http", Host: fmt.Sprintf("localhost:%d", pport)}
		cfg.APUrls, cfg.LPUrls = []url.URL{purl}, []url.URL{purl}

		cfg.AutoCompactionMode = compactor.ModePeriodic
		cfg.AutoCompactionRetention = 1 // every hour
		cfg.SnapCount = 1000            // single-node, keep minimum snapshot

		cfgs[i] = cfg
		iss[i] = fmt.Sprintf("%s=%s", cfg.Name, cfg.APUrls[0].String())
	}

	for i := range cfgs {
		cfgs[i].InitialCluster = strings.Join(iss, ",")
	}

	var wg sync.WaitGroup
	wg.Add(len(cfgs))

	for _, cfg := range cfgs {
		go func(cfg *embed.Config) {
			defer wg.Done()

			glog.Infof("starting %q with endpoint %q", cfg.Name, cfg.ACUrls[0].String())
			srv, err := embed.StartEtcd(cfg)
			if err != nil {
				panic(err)
			}
			select {
			case <-srv.Server.ReadyNotify():
				err = nil
			case err = <-srv.Err():
			case <-srv.Server.StopNotify():
				err = fmt.Errorf("received from etcdserver.Server.StopNotify")
			}
			if err != nil {
				panic(err)
			}
			glog.Infof("started %q with endpoint %q", cfg.Name, cfg.ACUrls[0].String())
		}(cfg)
	}
	wg.Wait()
}

func handleInterrupts(dataDir string) {
	notifier := make(chan os.Signal, 1)
	signal.Notify(notifier, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-notifier
		glog.Infof("received %q signal, shutting down...", sig)

		glog.Infof("removing %q", dataDir)
		os.RemoveAll(dataDir)

		signal.Stop(notifier)

		pid := syscall.Getpid()
		// exit directly if it is the "init" process, since the kernel will not help to kill pid 1.
		if pid == 1 {
			os.Exit(0)
		}
		syscall.Kill(pid, sig.(syscall.Signal))
	}()
}

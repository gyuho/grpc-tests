package main

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/coreos/etcd/auth/authpb"
	"github.com/coreos/etcd/clientv3"
	epb "github.com/coreos/etcd/etcdserver/api/v3election/v3electionpb"
	lockpb "github.com/coreos/etcd/etcdserver/api/v3lock/v3lockpb"
	pb "github.com/coreos/etcd/etcdserver/etcdserverpb"
	"github.com/coreos/etcd/integration"
	"github.com/coreos/etcd/pkg/testutil"

	"google.golang.org/grpc/grpclog"
)

func TestTest(t *testing.T) {
	clientv3.SetLogger(grpclog.NewLoggerV2WithVerbosity(os.Stderr, os.Stderr, os.Stderr, 4))

	defer testutil.AfterTest(t)

	clus := integration.NewClusterV3(t, &integration.ClusterConfig{Size: 1})
	defer clus.Terminate(t)

	authSetupRoot(t, toGRPC(clus.Client(0)).Auth)

	eps := clus.Client(0).Endpoints()
	cli, cerr := clientv3.New(clientv3.Config{Endpoints: eps, Username: "root", Password: "123"})
	if cerr != nil {
		t.Fatal(cerr)
	}
	defer cli.Close()

	fmt.Println("Sleep!!!")
	select {}
}

func authSetupUsers(t *testing.T, auth pb.AuthClient, users []user) {
	for _, user := range users {
		if _, err := auth.UserAdd(context.TODO(), &pb.AuthUserAddRequest{Name: user.name, Password: user.password}); err != nil {
			t.Fatal(err)
		}
		if _, err := auth.RoleAdd(context.TODO(), &pb.AuthRoleAddRequest{Name: user.role}); err != nil {
			t.Fatal(err)
		}
		if _, err := auth.UserGrantRole(context.TODO(), &pb.AuthUserGrantRoleRequest{User: user.name, Role: user.role}); err != nil {
			t.Fatal(err)
		}

		if len(user.key) == 0 {
			continue
		}

		perm := &authpb.Permission{
			PermType: authpb.READWRITE,
			Key:      []byte(user.key),
			RangeEnd: []byte(user.end),
		}
		if _, err := auth.RoleGrantPermission(context.TODO(), &pb.AuthRoleGrantPermissionRequest{Name: user.role, Perm: perm}); err != nil {
			t.Fatal(err)
		}
	}
}

func authSetupRoot(t *testing.T, auth pb.AuthClient) {
	root := []user{
		{
			name:     "root",
			password: "123",
			role:     "root",
			key:      "",
		},
	}
	authSetupUsers(t, auth, root)
	if _, err := auth.AuthEnable(context.TODO(), &pb.AuthEnableRequest{}); err != nil {
		t.Fatal(err)
	}
}

type user struct {
	name     string
	password string
	role     string
	key      string
	end      string
}

func toGRPC(c *clientv3.Client) grpcAPI {
	return grpcAPI{
		pb.NewClusterClient(c.ActiveConnection()),
		pb.NewKVClient(c.ActiveConnection()),
		pb.NewLeaseClient(c.ActiveConnection()),
		pb.NewWatchClient(c.ActiveConnection()),
		pb.NewMaintenanceClient(c.ActiveConnection()),
		pb.NewAuthClient(c.ActiveConnection()),
		lockpb.NewLockClient(c.ActiveConnection()),
		epb.NewElectionClient(c.ActiveConnection()),
	}
}

func newClientV3(cfg clientv3.Config) (*clientv3.Client, error) {
	return clientv3.New(cfg)
}

type grpcAPI struct {
	// Cluster is the cluster API for the client's connection.
	Cluster pb.ClusterClient
	// KV is the keyvalue API for the client's connection.
	KV pb.KVClient
	// Lease is the lease API for the client's connection.
	Lease pb.LeaseClient
	// Watch is the watch API for the client's connection.
	Watch pb.WatchClient
	// Maintenance is the maintenance API for the client's connection.
	Maintenance pb.MaintenanceClient
	// Auth is the authentication API for the client's connection.
	Auth pb.AuthClient
	// Lock is the lock API for the client's connection.
	Lock lockpb.LockClient
	// Election is the election API for the client's connection.
	Election epb.ElectionClient
}

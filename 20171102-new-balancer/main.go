package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/etcd/clientv3"
)

var eps = []string{
	"localhost:2379",
}

func main() {
	cli1, err := clientv3.New(clientv3.Config{Endpoints: eps})
	if err != nil {
		log.Fatalf("clientv3.New error: %v", err)
	}
	defer cli1.Close()

	/*
	 */
	_, err = cli1.Auth.UserAdd(context.Background(), "root", "123")
	if err != nil {
		log.Fatalf("UserAdd error: %v", err)
	}
	_, err = cli1.Auth.RoleAdd(context.Background(), "root")
	if err != nil {
		log.Fatalf("RoleAdd error: %v", err)
	}
	_, err = cli1.Auth.UserGrantRole(context.Background(), "root", "root")
	if err != nil {
		log.Fatalf("UserGrantRole error: %v", err)
	}
	_, err = cli1.Auth.AuthEnable(context.Background())
	if err != nil {
		log.Fatalf("AuthEnable error: %v", err)
	}
	fmt.Println("done!!!")

	/*
		time.Sleep(time.Second)
		acli1 := pb.NewAuthClient(cli1.ActiveConnection())
		_, err = acli1.UserAdd(context.Background(), &pb.AuthUserAddRequest{Name: "root", Password: "123"})
		if err != nil {
			log.Fatalf("UserAdd error: %v", err)
		}
		_, err = acli1.RoleAdd(context.Background(), &pb.AuthRoleAddRequest{Name: "root"})
		if err != nil {
			log.Fatalf("RoleAdd error: %v", err)
		}
		_, err = acli1.UserGrantRole(context.Background(), &pb.AuthUserGrantRoleRequest{User: "root", Role: "root"})
		if err != nil {
			log.Fatalf("UserGrantRole error: %v", err)
		}
		_, err = acli1.AuthEnable(context.Background(), &pb.AuthEnableRequest{})
		if err != nil {
			log.Fatalf("AuthEnable error: %v", err)
		}
		fmt.Println("done!!!")
	*/

	var cli2 *clientv3.Client
	cli2, err = clientv3.New(clientv3.Config{Endpoints: eps, Username: "root", Password: "123"})
	if err != nil {
		log.Fatalf("clientv3.New error: %v", err)
	}
	defer cli2.Close()

	fmt.Println("wait...")
	select {}
}

/*
2017-11-01 14:10:15.516423 N | auth: added a new user: root
2017-11-01 14:10:15.517235 N | auth: Role root is created
2017-11-01 14:10:15.523357 N | auth: granted role root to user root
2017-11-01 14:10:15.529770 N | auth: Authentication enabled
Server.Authenticate name:"root" password:"123"
Authenticate success? <nil> / name:"root" password:"123" simple_token:"DRVQWoacqFTfGCXZ"

handleSubConnStateChange: CONNECTING
--------- put CONNECTING
handleSubConnStateChange: READY
--------- put READY
============= auth.authenticate 1
============= auth.authenticate 2 <nil> <nil>
wait...
handleSubConnStateChange: CONNECTING
--------- put CONNECTING
handleSubConnStateChange: READY
--------- put READY
*/

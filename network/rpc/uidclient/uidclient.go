package uidclient

import (
	pb "go/cmkj_server_go/network/rpc/uidclient/uidgenerator"
	"go/cmkj_server_go/util"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "127.0.0.1:45454"
)

//NextUID 获取自增id
func NextUID(game int32) int64 {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		util.Log.Fatalf("grpc dial failed :%v", err)
	}
	defer conn.Close()

	client := pb.NewUIDGenneratorClient(conn)
	r, err := client.NextUid(context.Background(), &pb.Request{Game: game})
	if err != nil {
		util.Log.Error(err)
	}
	return r.Uid
}

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"rpcClient/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// hello_client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func runLotsOfReplies(c pb.GreeterClient) {
	// server端流式RPC
	//req := pb.HelloRequest{Name: *name}
	//ctx, cancel := context.WithTimeout(context.Background(), &req)
	//ctx,cancel:=context.Background()
	//defer cancel()
	stream, err := c.LotsOfReplies(context.Background(), &pb.HelloRequest{Name: "张三"})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v", err)
	}
	var times=0
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v", err)
		}
		times+=1
		if times==5{
			fmt.Println("客户端应该停止了")
			break
		}
		log.Printf("got reply: %q,times:%d\n", res.GetReply(),times)
	}
	stream.CloseSend()
}

func main() {
	flag.Parse()
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	runLotsOfReplies(c)
}
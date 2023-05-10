package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"net"
	"rpcTest/pb"
	"strings"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	fmt.Printf("\n %s服务器流 开始运行\n", in.Name)
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	//make chan
	//var wg sync.WaitGroup
	//wg.Add(1)
	iotserverChannel:=make(chan int)
	go runLotsOfReplies(c,iotserverChannel)

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	var times int = 0
	for range ticker.C {
		data := &pb.HelloResponse{
			Reply: in.GetName(),
		}
		times += 1
		fmt.Printf("服务器流发送数据,times:%d,name：%s\n", times, in.Name)
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			iotserverChannel<-1
			fmt.Printf("%s 服务器发生了错误：%+v\n", err)
			return err
		}
	}
	return nil
}

// BidiHello 双向流式打招呼
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		reply := magic(in.GetName()) // 对收到的数据做些处理

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

var addr = flag.String("addr", "127.0.0.1:9898", "the address to connect to")

//客户端流rpc
func runLotsOfReplies(c pb.GreeterClient,channel chan int ) {
	// 客户端端流式RPC
	//req := pb.HelloRequest{Name: *name}
	//ctx, cancel := context.WithTimeout(context.Background(), &req)
	//ctx,cancel:=context.Background()
	//defer cancel()
	stream, err := c.LotsOfReplies(context.Background(), &pb.HelloRequest{Name: "server端"})
	if err != nil {
		log.Fatalf("c.LotsOfReplies failed, err: %v", err)
	}
	var times = 0
	for {
		// 接收服务端返回的流式数据，当收到io.EOF或错误时退出
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("c.LotsOfReplies failed, err: %v", err)
		}
		_, ok := <-channel
		if ok {
			fmt.Println("通知iotserver关闭")
			stream.CloseSend()
		}
		log.Printf("got reply: %q,times:%d\n", res.GetReply(), times)
	}

}


func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	fmt.Println("-------------------------rpc 服务器启动-------------------------------")
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}

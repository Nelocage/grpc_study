package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	"rpcTest/pb"
	"time"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	fmt.Printf("\n %s服务器流 开始运行\n",in.Name)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()
	var times int =0
	for range ticker.C{
		data := &pb.HelloResponse{
			Reply:  in.GetName(),
		}
		times+=1
		fmt.Printf("iotserver发送数据,times:%d,name：%s\n",times,in.Name)
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			fmt.Printf("%s 服务器发生了错误：%+v\n",err)
			return err
		}
	}
	return nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":9898")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	fmt.Println("-------------------------rpc iotserver启动-------------------------------")
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}

}

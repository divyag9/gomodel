package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	pb "github.com/divyag9/gomodel/pkg/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

func main() {
	//Command line parameters
	serverAddr := flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	orderNumber := flag.Int64("order_number", 600051597, "Order number to retrieve image details for")
	flag.Parse()

	//Create grpc client
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewContentServiceClient(conn)

	//Add metadata to context
	header := metadata.New(map[string]string{"cache": "false"})
	ctx := context.Background()
	ctx = metadata.NewOutgoingContext(ctx, header)

	//Sending request
	response, err := client.ListByOrderNumber(ctx, &pb.OrderNumberRequest{OrderNumber: *orderNumber})
	if err != nil {
		log.Fatalf("error listing imageDetails by ordernumber: %v", err)
	}
	fmt.Println("response: ", response)
}

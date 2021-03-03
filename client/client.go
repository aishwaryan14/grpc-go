package main

import (
	"context"
	"fmt"
	"io"
	"log"

	computepb "github.com/aishwaryan14/grpc-go/protofiles"
	"google.golang.org/grpc"
)

//Function for Unary RPC call
func doUnaryRPC(c computepb.ComputeServiceClient) {
	fmt.Printf("Unary RPC call...\n")
	req := &computepb.EncryptionRequest{
		Pt: "This is a unary grpc call",
	}
	res, err := c.Encrypt(context.Background(), req)
	if err != nil {
		log.Fatalf("Couldn't connect to server")
	}
	fmt.Printf("Cipher text after encrytion:\n%v\n", res)
}

//Function for CLient Streaming for MinMaxSum computation
func doClientStreaming(c computepb.ComputeServiceClient) {
	fmt.Printf("Client Streaming RPC call for MinMaxSum...\n")
	stream, err := c.MinMaxSum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v\n", err)
	}
	numbers := []int32{1, 3, 5, 7, 9}

	for _, num := range numbers {
		fmt.Printf("Sending number %v\n", num)
		stream.Send(&computepb.MinMaxRequest{
			Number: num,
		})
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error receiving response %v\n", err)
	}
	fmt.Printf("The Minimum sum is %v\n", res.GetMinimum())
	fmt.Printf("The Maximum sum is %v\n", res.GetMaximum())

}

//Function for Server Streaming for Fibonacci series generation
func doServerStreaming(c computepb.ComputeServiceClient) {
	fmt.Printf("Server Streaming RPC for Fibonacci...")
	req := &computepb.FibonacciRequest{
		N: 10,
	}
	stream, err := c.Fibonacci(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Fibonacci RPC: %v", err)
	} else {
		fmt.Printf("Fibonacci Series:\n")
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error in receiving:%v\n", err)
			}
			fmt.Println(res.GetFib())
		}
	}
}

func main() {
	fmt.Println("Client..")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to the server %v \n", err)
	}
	defer cc.Close()
	c := computepb.NewComputeServiceClient(cc)
	doUnaryRPC(c)
	//doClientStreaming(c)
	//doServerStreaming(c)
}

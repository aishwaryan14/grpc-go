package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"sort"

	computepb "github.com/aishwaryan14/grpc-go/protofiles"
	"google.golang.org/grpc"
)

type server struct{}

//Unary RPC for Encryption
func (s *server) Encrypt(ctx context.Context, req *computepb.EncryptionRequest) (*computepb.EncryptionResponse, error) {
	fmt.Printf("Unary RPC request for Encryption:\n")
	ss := req.Pt
	li := []string{}
	for i := 0; i < len(ss); i++ {
		if ss[i] != ' ' {
			li = append(li, string(ss[i]))
		}
	}
	n := math.Sqrt(float64(len(li)))
	rows := int(math.Floor(n))
	cols := int(math.Ceil(n))
	matrix := make([][]string, cols)
	for i := range matrix {
		matrix[i] = make([]string, rows)
	}
	k := 0
	for i := 0; i < cols; i++ {
		for j := 0; j < rows; j++ {
			if k < len(li) {
				matrix[i][j] = li[k]
				k = k + 1
			}
		}
	}
	invmatrix := make([][]string, rows)
	for i := range invmatrix {
		invmatrix[i] = make([]string, cols)
	}
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			invmatrix[i][j] = matrix[j][i]
		}
	}

	str := ""
	for i := range invmatrix {
		for j := range invmatrix[i] {
			str = str + invmatrix[i][j]
		}
		str = str + " "
	}
	result := &computepb.EncryptionResponse{
		Ct: str,
	}
	return result, nil
}

//Client Streaming for MinMaxSum computation
func (s *server) MinMaxSum(stream computepb.ComputeService_MinMaxSumServer) error {
	fmt.Printf("MinMaxSum RPC called!\n")
	numSlice := []int{}
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			//MinMaxSumLogic
			sort.Ints(numSlice)
			li := []int{}
			sum1 := 0
			for i := range numSlice {
				sum1 = 0
				for j := range numSlice {
					if j != i {
						sum1 += numSlice[j]
					}
				}
				li = append(li, sum1)
			}
			sort.Ints(li)
			minimum := int32(li[0])
			maximum := int32(li[len(li)-1])
			return stream.SendAndClose(&computepb.MinMaxResponse{
				Minimum: minimum,
				Maximum: maximum,
			})
		}
		if err != nil {
			log.Fatalf("Error connecting to client %v\n", err)
		}
		numSlice = append(numSlice, int(req.Number))
	}
}

//Server Streaming for Fibonacci Series
func (s *server) Fibonacci(req *computepb.FibonacciRequest, stream computepb.ComputeService_FibonacciServer) error {
	fmt.Printf("Received Stream Server RPC call %v\n", req)
	number := req.GetN()
	fib := 0
	x := 1
	y := 0
	for i := 0; i < int(number); i++ {
		stream.Send(&computepb.FibonacciResponse{
			Fib: int32(fib),
		})
		fib = x + y
		x = y
		y = fib
	}
	return nil
}

func main() {
	fmt.Println("Server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen %v\n", err)
	}
	s := grpc.NewServer()

	computepb.RegisterComputeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve:%v\n", err)
	}
}

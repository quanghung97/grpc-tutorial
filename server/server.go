package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"quanghung97/.calculator/calculatorpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

type server struct {
	calculatorpb.CalculatorServiceServer
}

// unary API type
// implement interface
// con trỏ server khởi tạo mới, mỗi lần request thì sẽ tạo vùng nhớ khác nhau
func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum called...")
	// nội dung của response đều cố định được tham chiếu bằng địa chỉ
	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

// sever streaming API
func (*server) PrimeNumberDecomposition(
	req *calculatorpb.PNDRequest,
	stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {

	log.Println("PrimeNumberDecomposition called...")
	k := int32(2)
	N := req.GetNumber()
	for N > 1 {
		if N%k == 0 {
			N = N / k

			// send response to client
			stream.SendMsg(&calculatorpb.PNDResponse{
				Result: k,
			})
			time.Sleep(500 * time.Millisecond)
		} else {
			k++
			log.Println("k increase to %v", k)
		}
	}

	return nil
}

// client streaming API (client gui nhieu req len server)
func (*server) Average(stream calculatorpb.CalculatorService_AverageServer) error {
	log.Println("Average called...")
	var total float32
	var count int
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// tinh trung binh va return cho client
			resp := &calculatorpb.AverageResponse{
				Result: total / float32(count),
			}

			return stream.SendAndClose(resp)
		}

		if err != nil {
			log.Fatalf("err while Recv Average %v", err)
			return err
		}

		log.Println("receive num %v", req)
		total += req.GetNum()
		count++
	}
}

// stream API both server and client (bi directional streaming API)
func (*server) FindMax(stream calculatorpb.CalculatorService_FindMaxServer) error {
	log.Println("find max called ...")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF... client stop call")
			return nil
		}

		if err != nil {
			log.Fatalf("err while Recv FindMax %v", err)
			return err
		}

		num := req.GetNum()
		log.Printf("recv num %v\n", num)
		if num > max {
			max = num
		}

		// send response to client
		err = stream.SendMsg(&calculatorpb.FindMaxResponse{
			Max: max,
		})

		if err != nil {
			log.Fatalf("send max err %v", err)
			return err
		}
	}
}

// handle error
func (*server) Square(ctx context.Context, req *calculatorpb.SquareRequest) (*calculatorpb.SquareResponse, error) {
	log.Printf("square called ...")

	num := req.GetNum()

	if num < 0 {
		log.Print("req num < 0, num=%v, return invalidArgument", num)
		return nil, status.Errorf(codes.InvalidArgument, "Expect num > 0, req num was %v", num)
	}

	resp := &calculatorpb.SquareResponse{
		SquareRoot: math.Sqrt(float64(num)),
	}

	return resp, nil
}

func (*server) SumWithDeadline(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	log.Println("sum with deadline called...")

	// fake time waiting chờ 3s mới response cho client
	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			log.Println("context.Canceled...")
			return nil, status.Error(codes.Canceled, "client canceled req")
		}
		time.Sleep(1 * time.Second)
	}

	resp := &calculatorpb.SumResponse{
		Result: req.GetNum1() + req.GetNum2(),
	}

	return resp, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:5000")

	if err != nil {
		log.Fatalf("err listen %v", err)
	}

	certFile := "ssl/server.crt"
	keyFile := "ssl/server.key"

	creds, sslErr := credentials.NewServerTLSFromFile(certFile, keyFile)
	if sslErr != nil {
		log.Fatalf("create creds ssl err %v\n", sslErr)
		return
	}
	opts := grpc.Creds(creds)

	s := grpc.NewServer(opts)

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	fmt.Println("calculator is running")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err serve %v", err)
	}
}

package main

import (
	"context"
	"io"
	"log"
	"quanghung97/.calculator/calculatorpb"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	cc, err := grpc.Dial("localhost:5000", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v", err)
	}

	// same finally in javascript, close connection
	defer cc.Close()

	client := calculatorpb.NewCalculatorServiceClient(cc)

	// log.Fatalf("service client %f", client)

	// callSum(client)
	// callPND(client)
	// callAverage(client)
	// callFindMax(client)
	callSquareRoot(client, -9)
}

func callSum(c calculatorpb.CalculatorServiceClient) {
	log.Println("calling sum api")
	resp, err := c.Sum(context.Background(), &calculatorpb.SumRequest{
		Num1: 7,
		Num2: 6,
	})

	if err != nil {
		log.Fatalf("call sum api err %v", err)
	}

	log.Printf("sum api response %v", resp.GetResult())
}

// server streaming API
func callPND(c calculatorpb.CalculatorServiceClient) {
	log.Println("pnd api")
	stream, err := c.PrimeNumberDecomposition(context.Background(), &calculatorpb.PNDRequest{
		Number: 120,
	})

	if err != nil {
		log.Fatalf("callPND err %v", err)
	}

	for {
		resp, recvErr := stream.Recv()
		if recvErr == io.EOF {
			log.Println("server finish streaming")
			return
		}

		log.Printf("Prime number %v", resp.GetResult())
	}
}

// client streaming API
func callAverage(c calculatorpb.CalculatorServiceClient) {
	log.Println("callAverage api")
	stream, err := c.Average(context.Background())

	if err != nil {
		log.Fatalf("call averrage err %v", err)
	}

	listReq := []calculatorpb.AverageRequest{
		{
			Num: 5,
		},
		calculatorpb.AverageRequest{
			Num: 10,
		},
		calculatorpb.AverageRequest{
			Num: 12,
		},
		calculatorpb.AverageRequest{
			Num: 3,
		},
		calculatorpb.AverageRequest{
			Num: 4.2,
		},
	}

	for _, req := range listReq {
		err := stream.Send(&req)

		if err != nil {
			log.Fatalf("send average request err %v", err)
			break
		}

		time.Sleep(500 * time.Millisecond)
	}

	resp, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("receive average response err %v", err)
	}

	log.Printf("average response %+v", resp)
}

// stream API both server and client (bi directional streaming API)
func callFindMax(c calculatorpb.CalculatorServiceClient) {
	log.Println("callFindMax api")
	stream, err := c.FindMax(context.Background())

	if err != nil {
		log.Fatalf("call findmax err %v", err)
	}

	//go waiting để call lại go func
	waitc := make(chan struct{})

	go func() {
		// gui nhieu request
		listReq := []calculatorpb.FindMaxRequest{
			{
				Num: 5,
			},
			calculatorpb.FindMaxRequest{
				Num: 10,
			},
			calculatorpb.FindMaxRequest{
				Num: 12,
			},
			calculatorpb.FindMaxRequest{
				Num: 3,
			},
			calculatorpb.FindMaxRequest{
				Num: 4,
			},
		}
		for _, req := range listReq {
			err := stream.Send(&req)

			if err != nil {
				log.Fatalf("send find max request err %v", err)
				break
			}
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				log.Println("ending find max api ...")
				break
			}
			if err != nil {
				log.Fatalf("recv find max err %v", err)
				break
			}

			log.Printf("max: %v\n", resp.GetMax())
		}
		close(waitc)
	}()

	<-waitc
}

func callSquareRoot(c calculatorpb.CalculatorServiceClient, num int32) {
	log.Println("calling square api")
	resp, err := c.Square(context.Background(), &calculatorpb.SquareRequest{
		Num: num,
	})

	if err != nil {
		log.Printf("call square root api err %v\n", err)
		errStatus, ok := status.FromError(err)
		if ok {
			log.Printf("err msg: %v\n", errStatus.Message())
			log.Printf("err code: %v\n", errStatus.Code())
			if errStatus.Code() == codes.InvalidArgument {
				log.Printf("InvalidArgument num %v", num)
				return
			}
		}
	}

	log.Printf("square api response %v", resp.GetSquareRoot())
}

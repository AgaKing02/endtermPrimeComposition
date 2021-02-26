package main

import (
	"context"
	"endtermPrimeComposition/calculatorpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	fmt.Println("Client is on")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculateServiceClient(conn)
	PrimeComposition(c)
}

func PrimeComposition(c calculatorpb.CalculateServiceClient) {
	ctx := context.Background()
	req := &calculatorpb.CalculatorRequest{Calculating: &calculatorpb.Calculating{
		Number: 15,
	}}

	stream, err := c.PrimeCompose(ctx, req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC %v", err)
	}
	defer stream.CloseSend()

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break LOOP
			}
			log.Fatalf("error while reciving from GreetManyTimes RPC %v", err)
		}
		log.Printf("response from GreetManyTimes:%v \n", res.GetResult())
	}

}

func getAverage(c calculatorpb.CalculateServiceClient) {

	requests := []*calculatorpb.NumberRequest{
		{
			Numbers: 2,
		},
		{
			Numbers: 3,
		},
		{
			Numbers: 12,
		},
		{
			Numbers: 9,
		},
	}

	ctx := context.Background()
	stream, err := c.ComputeAverage(ctx)
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage Response: %v\n", res)
}

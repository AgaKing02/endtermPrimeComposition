package main

import (
	"endtermPrimeComposition/calculatorpb"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"time"
)
//Functions
func Average(arr []int64) (avg float64) {
	sm := 0.0
	ln := len(arr)
	for i := 0; i < ln; i++ {
		sm += float64(arr[i])
	}
	avg = sm / float64(ln)
	return
}
func Factors(n int) (pfs []int) {
	for n%2 == 0 {
		pfs = append(pfs, 2)
		n = n / 2
	}
	for i := 3; i*i <= n; i = i + 2 {
		// while i divides n, append i and divide n
		for n%i == 0 {
			pfs = append(pfs, i)
			n = n / i
		}
	}
	if n > 2 {
		pfs = append(pfs, n)
	}
	return
}


type Server struct {
	calculatorpb.UnimplementedCalculateServiceServer
}
//Structure
func (s *Server) GetAverage(req *calculatorpb.CalculatorRequest, stream calculatorpb.CalculateService_PrimeComposeServer) error {
	fmt.Printf("GreetManyTimes function was invoked with %v \n", req)
	number := req.GetCalculating().GetNumber()
	n := int(number)
	arr := Factors(n)
	for i := 0; i < len(arr); i++ {
		res := &calculatorpb.CalculatorResponse{Result: fmt.Sprintf("%d) Hello, %v\n", n, arr[i])}
		if err := stream.Send(res); err != nil {
			log.Fatalf("error while sending greet many times responses: %v", err.Error())
		}
		time.Sleep(time.Second)
	}
	return nil
}



func (s *Server) ComputeAverage(stream calculatorpb.CalculateService_ComputeAverageServer) error {
	fmt.Printf("AverageClient function was invoked with a streaming request\n")
	var result float64
	var arr []int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// we have finished reading the client stream
			result = Average(arr)
			return stream.SendAndClose(&calculatorpb.AverageResponse{Result: result})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}
		arr = append(arr, req.GetNumbers())
	}
}


func main() {
	l, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	calculatorpb.RegisterCalculateServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}
}

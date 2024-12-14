package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	ecpb "google.golang.org/grpc/examples/features/proto/echo"
)

var port = flag.Int("port", 50051, "the port to serve on")

type ecServer struct {
	ecpb.UnimplementedEchoServer
}

// UnaryEcho implements the unary echo method
func (s *ecServer) UnaryEcho(ctx context.Context, req *ecpb.EchoRequest) (*ecpb.EchoResponse, error) {
	return &ecpb.EchoResponse{Message: req.Message}, nil
}

// ServerStreamingEcho implements the server streaming echo method
func (s *ecServer) ServerStreamingEcho(req *ecpb.EchoRequest, stream ecpb.Echo_ServerStreamingEchoServer) error {
	// Send messages for 10 seconds
	for i := 0; i < 10; i++ {
		message := fmt.Sprintf("%s - stream response %d", req.Message, i+1)
		if err := stream.Send(&ecpb.EchoResponse{Message: message}); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

// ClientStreamingEcho implements the client streaming echo method
func (s *ecServer) ClientStreamingEcho(stream ecpb.Echo_ClientStreamingEchoServer) error {
	var messages []string
	startTime := time.Now()

	for {
		// Exit after 10 seconds
		if time.Since(startTime) > 10*time.Second {
			break
		}

		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("Received messages: %v", messages)
			return stream.SendAndClose(&ecpb.EchoResponse{Message: message})
		}
		if err != nil {
			return err
		}
		messages = append(messages, req.Message)
	}

	message := fmt.Sprintf("Timeout reached. Received messages: %v", messages)
	return stream.SendAndClose(&ecpb.EchoResponse{Message: message})
}

// BidirectionalStreamingEcho implements the bidirectional streaming echo method
func (s *ecServer) BidirectionalStreamingEcho(stream ecpb.Echo_BidirectionalStreamingEchoServer) error {
	startTime := time.Now()

	for {
		// Exit after 10 seconds
		if time.Since(startTime) > 10*time.Second {
			return nil
		}

		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		message := fmt.Sprintf("Bidirectional: %s", req.Message)
		if err := stream.Send(&ecpb.EchoResponse{Message: message}); err != nil {
			return err
		}
	}
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	fmt.Printf("server listening at %v\n", lis.Addr())

	s := grpc.NewServer()

	// Register Echo server
	ecpb.RegisterEchoServer(s, &ecServer{})

	// Register reflection service on gRPC server
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net"

	// Import proto package for serialization

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"

	// Import your generated Go gRPC files
	pb "auth_hmac/pb"
)

// Your shared secret
var secretKey = []byte("root")

// Server implements the ExampleServiceServer interface
type Server struct {
	pb.UnimplementedExampleServiceServer
}

func (s *Server) SayHello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	// Extract HMAC from metadata
	fmt.Println("Received request")
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}
	hmacDigests := md.Get("hmac")
	if len(hmacDigests) == 0 {
		return nil, fmt.Errorf("HMAC not provided")
	}
	clientHMAC := hmacDigests[0]

	// Serialize the request
	serializedRequest, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	// Generate HMAC for the received request
	mac := hmac.New(sha256.New, secretKey)
	mac.Write(serializedRequest)
	expectedHMAC := fmt.Sprintf("%x", mac.Sum(nil))

	// Verify HMAC
	fmt.Println("Client HMAC:", clientHMAC)
	fmt.Println("Expected HMAC:", expectedHMAC)
	if !hmac.Equal([]byte(clientHMAC), []byte(expectedHMAC)) {
		return nil, fmt.Errorf("invalid HMAC")
	}

	// Process the request
	// Your business logic goes here...
	fmt.Print("Received request\n")
	return &pb.HelloResponse{Message: "root"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, &Server{})
	fmt.Println("gRPC server started")
	if err := grpcServer.Serve(lis); err != nil {
		fmt.Printf("failed to serve: %v", err)
	}

}

// main.go
package main

import (
	"context"
	"fmt"
	pb "grpc-server/pb"
	"log"
	"net"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}
type BookingServer struct {
	pb.UnimplementedBookingServiceServer
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	return &pb.UserResponse{
		UserId: req.UserId,
		Name:   "John Doe",
	}, nil
}

func (s *BookingServer) Booking(ctx context.Context, req *pb.BookingRequest) (*pb.BookingResponse, error) {
	// Connect to RabbitMQ server
	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	// Create a channel
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare a queue named "hello"
	q, err := ch.QueueDeclare(
		"bookings", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Publish a message to the queue
	// data, err := proto.Unmarshal()
	// if err != nil {
	// 	log.Fatalf("Failed to marshal proto : %v", err)
	// }
	data, err := protojson.Marshal(req)
	if err != nil {
		log.Fatalf("Failed to marshal proto : %v", err)
	}
	err = ch.PublishWithContext(
		ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(data),
		})
	if err != nil {
		log.Fatalf("Failed to publish a message: %v", err)
	}

	fmt.Println(" [x] Sent new booking request", req.String())

	return &pb.BookingResponse{
		BookingId:   req.TicketId,
		TicketId:    123,
		BookingName: req.BookingName,
		Email:       req.Email,
	}, err
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServer{})
	pb.RegisterBookingServiceServer(s, &BookingServer{})
	log.Println("Server running on port 50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

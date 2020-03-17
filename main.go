// shippy-service-vessel/main.go
package main

import (
	"context"
	"errors"
	pb "github.com/suliar/shippy-service-vessel/proto/vessel"
	"google.golang.org/grpc"
	"log"
	"net"
)

type Repository interface {
	FindAvailable(*pb.Specification) (*pb.Vessel, error)
}

type VesselRepository struct {
	vessels []*pb.Vessel
}


//FindAvailable - checks a specification against a map vessels,
// if capacity and max weight are below a vessels capacity and max weight,
//then return that vessel

func (repo *VesselRepository) FindAvailable(spec *pb.Specification) (*pb.Vessel, error) {
	for _, v := range repo.vessels {
		if spec.Capacity <= v.Capacity && spec.MaxWeight <= v.MaxWeight {
			return v, nil
		}
	}
	return nil, errors.New("No vessel found by that spec")
}

//Our grpc service handler
type service struct {
	repo Repository
}

func(s service) FindAvailable(ctx context.Context, req *pb.Specification) (*pb.Response, error) {

	// Find the next available vessel
	vessel, err := s.repo.FindAvailable(req)
	if err != nil {
		return nil, err
	}

	return &pb.Response{Vessel:
		vessel}, nil
}


func main() {

	vessels := []*pb.Vessel{
		{
			Id:        "vessel001",
			Capacity:  500,
			MaxWeight: 200000,
			Name:      "Boaty McBoatface",
		},
	}

	repo := &VesselRepository{vessels:
		vessels}

	port := ":50052"

	//conn, err := grpc.Dial(port, grpc.WithInsecure())
	//if err != nil {
	//	log.Fatalf("Did not connect: %v", err)
	//}
	// Set-up our gRPC server.
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	// Register our service with the gRPC server, this will tie our
	// implementation into the auto-generated interface code for our
	// protobuf definition.
	pb.RegisterVesselServiceServer(s, &service{repo})

	log.Println("Running on port:", port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
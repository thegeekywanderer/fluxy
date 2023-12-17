// Package grpc defines the GRPC handler for fluxy
package grpc

import (
	"context"
	"errors"
	"strconv"

	"github.com/thegeekywanderer/fluxy/models"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
	"github.com/thegeekywanderer/fluxy/proto"
	"google.golang.org/grpc"
)

// FluxyServeStruct defines the grpc server struct for fluxy
type FluxyServeStruct struct {
	useCase interfaces.UseCaseInterface
	proto.UnimplementedRateLimiterServiceServer
  }
 
// NewServer function registers a new ratelimiters service to the grpc server
func NewServer(grpcServer *grpc.Server, usecase interfaces.UseCaseInterface){
	userGrpc := &FluxyServeStruct{useCase: usecase}
	proto.RegisterRateLimiterServiceServer(grpcServer, userGrpc)
}

// RegisterClient function implements the handler for registering a new client to fluxy
func (s *FluxyServeStruct) RegisterClient(_ context.Context, req *proto.ClientRequest) (*proto.ClientResponse, error) {
	data := s.transformClientRPC(req)
	if data.Name == "" {
		return &proto.ClientResponse{}, errors.New("name cannot be blank")
	}
	client, err := s.useCase.RegisterClient(data)
	if err != nil {
		return &proto.ClientResponse{}, err
	}
	return s.transformClientModel(client), nil
}
 
// GetClient function returns the client instance registered into fluxy
func (s *FluxyServeStruct) GetClient(_ context.Context, req *proto.SingleClientRequest) (*proto.ClientResponse, error) {
	name := req.GetName()
	if name == "" {
		return &proto.ClientResponse{}, errors.New("name cannot be blank")
	}
	client, err := s.useCase.GetClient(name)
	if err != nil {
		return &proto.ClientResponse{}, err
	}
	return s.transformClientModel(client), nil
}

// UpdateClient function implements the handler for registering a new client to fluxy
func (s *FluxyServeStruct) UpdateClient(_ context.Context, req *proto.ClientRequest) (*proto.SuccessResponse, error) {
	data := s.transformClientRPC(req)
	err := s.useCase.UpdateClient(data)
	if err != nil {
		return &proto.SuccessResponse{}, err
	}
	return &proto.SuccessResponse{Response: "Updated client successfully"}, nil
}

// DeleteClient function returns the client instance registered into fluxy
func (s *FluxyServeStruct) DeleteClient(_ context.Context, req *proto.SingleClientRequest) (*proto.SuccessResponse, error) {
	name := req.GetName()
	if name == "" {
		return &proto.SuccessResponse{}, errors.New("name cannot be blank")
	}
	err := s.useCase.DeleteClient(name)
	if err != nil {
		return &proto.SuccessResponse{}, err
	}
	return &proto.SuccessResponse{Response: "Deleted client successfully"}, nil
}

// VerifyLimit function checks if the client is within the rate limit or not
func (s *FluxyServeStruct) VerifyLimit(ctx context.Context, req *proto.SingleClientRequest) (*proto.StateResponse, error) {
	name := req.GetName()
	if name == "" {
		return &proto.StateResponse{}, errors.New("name cannot be blank")
	}
	res, err := s.useCase.VerifyLimit(name)
	if err != nil {
		return &proto.StateResponse{}, err
	}
	allowed := false
	if res.State == interfaces.Allow {
		allowed = true
	}

	return &proto.StateResponse{
		Allowed: allowed, 
		TotalRequests: int64(res.TotalRequests), 
		ExpiresAt: res.ExpiresAt.Unix(),
	}, nil
}

func (s *FluxyServeStruct) transformClientRPC(req *proto.ClientRequest) models.Client{
	return models.Client{
		Name: req.GetName(), 
		Limit: uint64(req.GetLimit()), 
		Duration: uint64(req.GetDuration()),
	}
}

func (s *FluxyServeStruct) transformClientModel(client models.Client) *proto.ClientResponse{
	return &proto.ClientResponse{
		Id: strconv.Itoa(int(client.ID)), 
		Name: client.Name, 
		Limit: int64(client.Limit), 
		Duration: int64(client.Duration),
	}
}
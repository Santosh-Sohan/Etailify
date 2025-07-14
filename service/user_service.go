package service

import (
	"context"

	user "github.com/Santosh-Sohan/user-service/api"
	"github.com/Santosh-Sohan/user-service/repository"
)

type UserServiceServer struct {
	user.UnimplementedUserServiceServer
	Repo repository.UserRepository
}

func NewUserService() *UserServiceServer {
	return &UserServiceServer{
		Repo: repository.NewUserRepository(),
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *user.CreateUserRequest) (*user.UserResponse, error) {
	err := s.Repo.CreateUser(req.User)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{User: req.User}, nil
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, req *user.UpdateUserRequest) (*user.UserResponse, error) {
	err := s.Repo.UpdateUser(req)
	if err != nil {
		return nil, err
	}
	// optional: fetch user again
	return &user.UserResponse{User: &user.User{Id: req.Id}}, nil
}

func (s *UserServiceServer) UpdateContact(ctx context.Context, req *user.UpdateContactRequest) (*user.UserResponse, error) {
	err := s.Repo.UpdateContact(req)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{User: &user.User{Id: req.Id}}, nil
}

func (s *UserServiceServer) BlockUser(ctx context.Context, req *user.UserIdRequest) (*user.UserResponse, error) {
	err := s.Repo.BlockUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{User: &user.User{Id: req.Id}}, nil
}

func (s *UserServiceServer) UnblockUser(ctx context.Context, req *user.UserIdRequest) (*user.UserResponse, error) {
	err := s.Repo.UnblockUser(req.Id)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{User: &user.User{Id: req.Id}}, nil
}

func (s *UserServiceServer) GetUserByEmailOrPhone(ctx context.Context, req *user.EmailOrPhoneRequest) (*user.UserResponse, error) {
	u, err := s.Repo.GetUserByEmailOrPhone(req.PhoneNumber, req.Email)
	if err != nil {
		return nil, err
	}
	return &user.UserResponse{User: u}, nil
}

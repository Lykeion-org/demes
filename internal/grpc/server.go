package grpc

import (
	"net"
	"context"
	lyk_grpc "github.com/Lykeion-org/go-shared/pkg/grpc"
	pb "github.com/Lykeion-org/go-shared/pkg/grpc/generated/demes"
	repo "github.com/Lykeion-org/demes/internal/repository"
	"gorm.io/gorm"
	"google.golang.org/grpc"
)

type DemesServer struct {
	pb.UnimplementedDemesServiceServer
	UserRepository repo.UserRepository
	activeServer *grpc.Server
}

func NewDemesServer(db *gorm.DB) lyk_grpc.Server{
	return &DemesServer{
		UserRepository: repo.NewUserRepository(db),
	}
}

func (s *DemesServer)StartServer(target string) error {
	grpcServer := grpc.NewServer()
	pb.RegisterDemesServiceServer(grpcServer,s)

	listener, err := net.Listen("tcp", target)
	if err != nil {
		return err
	}

	go func(){
		err := grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
	
	s.activeServer = grpcServer
	return nil
}

func (s *DemesServer) StopServer() error {
	s.activeServer.GracefulStop()
	return nil
}

func (s *DemesServer) ReadUser(ctx context.Context, req *pb.UserUidRequest) (*pb.User, error) {
	user, err := s.UserRepository.ReadUser(ctx, req.Uid)
	if err != nil {
		return nil, err
	}

	return ConvertUserToProto(user), nil
}

func (s *DemesServer) ReadUserByEmailAddress(ctx context.Context, req *pb.UserMailRequest)  (*pb.User, error) {
	user, err := s.UserRepository.ReadUserByEmailAddress(ctx, req.EmailAddress)
	if err != nil {
		return nil, err
	}

	return ConvertUserToProto(user), nil
}

func (s *DemesServer) UpdateUser(ctx context.Context, req *pb.User) (*pb.User, error) {
	userModel, err := ConvertProtoToUser(req)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.UserRepository.UpdateUser(ctx, userModel)
	if err != nil {
		return nil, err
	}

	return ConvertUserToProto(updatedUser), nil
}

func (s *DemesServer) DeleteUser(ctx context.Context, req *pb.UserUidRequest) (*pb.User, error) {
	deletedUser, err := s.UserRepository.DeleteUser(ctx, req.Uid)
	if err != nil {
		return nil, err
	}
	return ConvertUserToProto(deletedUser), nil
}

func (s *DemesServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	page := int(req.Offset/req.Limit) + 1
	if req.Limit == 0 {
		req.Limit = 10
	}

	users, err := s.UserRepository.ListUsers(ctx, page, int(req.Limit))
	if err != nil {
		return nil, err
	}

	var pbUsers []*pb.User = ConvertUserListToProto(users)

	return &pb.ListUsersResponse{
		Users:       pbUsers,
		TotalCount:  int32(len(pbUsers)),
		PageSize:    req.Limit,
		CurrentPage: int32(page),
	}, nil
}

func (s *DemesServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	newUser,err := s.UserRepository.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return ConvertUserToProto(newUser), nil
}

func (s *DemesServer) ValidatePassword(ctx context.Context, req *pb.ValidatePasswordRequest) (*pb.ValidatePasswordResponse, error) {
	res, err := s.UserRepository.CheckPassword(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return &pb.ValidatePasswordResponse{
		Correct: res,
	}, nil

}

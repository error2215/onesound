package grpc

import (
	"context"

	"golang.org/x/crypto/bcrypt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"onesound/compiled/common"
	"onesound/compiled/onesound_api"

	"onesound/server/elastic/room"
	"onesound/server/elastic/user"
)

func (g *GRPCServer) HealthCheck(ctx context.Context, in *common.HealthCheckRequest) (*common.HealthCheckResponse, error) {
	logrus.Printf("Health check for %s", in.Service)
	return &common.HealthCheckResponse{
		Status: common.HealthCheckResponse_SERVING,
	}, nil
}

func (g *GRPCServer) CreateRoom(ctx context.Context, in *onesound_api.CreateRoomRequest) (*onesound_api.CreateRoomResponse, error) {
	exist, err := room.New().Names([]string{in.Name}).CheckIfRoomExist(ctx)
	if err != nil {
		return nil, err
	}
	if !exist {
		if in.Name == "" {
			return nil, errors.New("Name must not be empty")
		}
		resp, err := room.New().Password(in.Pass).Names([]string{in.Name}).CreateRoom(ctx)
		if err != nil {
			return nil, err
		}
		return &onesound_api.CreateRoomResponse{Room: resp}, nil
	}
	return nil, errors.New("Room with such name is already exist. Please try another one")
}

func (g *GRPCServer) DeleteRoom(ctx context.Context, in *onesound_api.DeleteRoomRequest) (*common.SimpleResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) AddVideo(ctx context.Context, in *onesound_api.AddVideoRequest) (*common.SimpleResponse, error) {

}

func (g *GRPCServer) RemoveVideo(ctx context.Context, in *onesound_api.SkipVideoRequest) (*common.SimpleResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) GetPlaylist(ctx context.Context, in *onesound_api.GetPlaylistRequest) (*onesound_api.GetPlaylistResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) CreateVoting(ctx context.Context, in *onesound_api.CreateVotingRequest) (*common.SimpleResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) AddPointToVoting(ctx context.Context, in *onesound_api.AddPointToVotingRequest) (*onesound_api.AddPointToVotingResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) Register(ctx context.Context, in *onesound_api.RegisterRequest) (*common.SimpleResponse, error) {
	exist, err := user.New().Emails([]string{in.Email}).CheckIfUserExist(ctx)
	if err != nil {
		return &common.SimpleResponse{OK: false}, err
	}
	if !exist {
		created, err := user.New().Names([]string{in.Name}).Emails([]string{in.Email}).Password(in.Password).CreateUser(ctx)
		if err != nil {
			return &common.SimpleResponse{OK: false}, err
		}
		if created {
			return &common.SimpleResponse{OK: true}, nil
		}
	}
	return &common.SimpleResponse{OK: false}, errors.New("User with such email already exist")
}

func (g *GRPCServer) Auth(ctx context.Context, in *onesound_api.AuthRequest) (*common.SimpleResponse, error) {
	users, err := user.New().Emails([]string{in.Email}).FindUsers(ctx)
	if err != nil {
		return &common.SimpleResponse{OK: false}, err
	}
	if len(users) > 0 {
		err = bcrypt.CompareHashAndPassword([]byte(users[0].Password), []byte(in.Password))
		if err != nil && err != bcrypt.ErrMismatchedHashAndPassword {
			return &common.SimpleResponse{OK: false}, err
		}
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return &common.SimpleResponse{OK: false}, errors.New("Password is incorrect")
		}
		return &common.SimpleResponse{OK: true}, nil
	}
	return &common.SimpleResponse{OK: false}, errors.New("User with such email do not exist")
}

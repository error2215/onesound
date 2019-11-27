package grpc

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"onesound/compiled/common"
	"onesound/compiled/onesound_api"
	"onesound/server/elastic/room"
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

func (g *GRPCServer) AddVideo(context.Context, *onesound_api.AddVideoRequest) (*common.SimpleResponse, error) {

}

func (g *GRPCServer) RemoveVideo(context.Context, *onesound_api.SkipVideoRequest) (*common.SimpleResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) GetPlaylist(context.Context, *onesound_api.GetPlaylistRequest) (*onesound_api.GetPlaylistResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) CreateVoting(context.Context, *onesound_api.CreateVotingRequest) (*common.SimpleResponse, error) {
	panic("implement me")
}

func (g *GRPCServer) AddPointToVoting(context.Context, *onesound_api.AddPointToVotingRequest) (*onesound_api.AddPointToVotingResponse, error) {
	panic("implement me")
}

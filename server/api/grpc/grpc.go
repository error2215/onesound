package grpc

import (
	"context"

	"github.com/sirupsen/logrus"

	"onesound/compiled/common"
)

func (g *GRPCServer) HealthCheck(context.Context, *common.HealthCheckRequest) (*common.HealthCheckResponse, error) {
	logrus.Info("kek")
	return &common.HealthCheckResponse{
		Status: 0,
	}, nil
}

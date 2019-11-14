package grpc

import (
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"onesound/compiled/onesound_api"
)

var GlobalGRPCServer *GRPCServer

type GRPCServer struct {
}

// Serve starts server
func (g *GRPCServer) Start(port string) error {

	GlobalGRPCServer = g
	addr := ":" + port
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.WithError(err).Fatalf("Cannot bind to %s", addr)
		return err
	}

	log.WithFields(log.Fields{
		"addr": addr,
	}).Info("GRPCServer Start")

	grpcServer := grpc.NewServer(
		grpc.MaxConcurrentStreams(20),
	)

	onesound_api.RegisterOneSoundAPIServer(grpcServer, g)
	_ = grpcServer.Serve(lis)

	return nil
}

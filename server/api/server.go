package api

import (
	"sync"

	log "github.com/sirupsen/logrus"

	"onesound/server/api/grpc"
	"onesound/server/config"
)

func Start() {
	grpcPort := config.GlobalConfig.AppPort

	log.WithFields(log.Fields{
		"grpcPort": grpcPort,
	}).Info("Launching GRPC server")

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		g := &grpc.GRPCServer{}
		g.Start(grpcPort)
	}()

	wg.Wait()
}

package server

import (
	"context"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health/grpc_health_v1"
)

type healthServer struct {
	logEntry *logrus.Entry
}

func NewHealthServer(logEntry *logrus.Entry) (grpc_health_v1.HealthServer, error) {
	healthServerObj := &healthServer{
		logEntry: logEntry,
	}
	return healthServerObj, nil
}

func (h *healthServer) Check(ctx context.Context, in *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	h.logEntry.Infof("Serving the Check request for health check: %+v", in)
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (h *healthServer) Watch(in *grpc_health_v1.HealthCheckRequest, _ grpc_health_v1.Health_WatchServer) error {
	h.logEntry.Infof("Serving the Check request for health check: %+v", in)
	return nil
}

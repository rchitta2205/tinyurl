package interceptor

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"time"
	"tinyurl/pkg/datamodel"
)

// AuthInterceptor is a server interceptor for authorization
type AuthInterceptor struct {
	authApp  datamodel.AuthApplication
	logEntry *logrus.Entry
}

// NewInterceptor returns a new auth interceptor
func NewAuthInterceptor(authApp datamodel.AuthApplication, logEntry *logrus.Entry) *AuthInterceptor {
	return &AuthInterceptor{
		authApp:  authApp,
		logEntry: logEntry,
	}
}

func (auth *AuthInterceptor) UnaryAuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := auth.authorize(ctx, info.FullMethod); err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(ctx, req)
}

func (auth *AuthInterceptor) StreamAuthInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if err := auth.authorize(stream.Context(), info.FullMethod); err != nil {
		return status.Error(codes.PermissionDenied, err.Error())
	}
	return handler(srv, stream)
}

func (auth *AuthInterceptor) authorize(ctx context.Context, method string) error {
	auth.logEntry.Info("Authenticating and authorizing user...")
	peerObj, ok := peer.FromContext(ctx)
	if !ok {
		return errors.New("error to read peer information")
	}

	tlsInfo, ok := peerObj.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return errors.New("error to get auth information")
	}

	certs := tlsInfo.State.VerifiedChains
	if len(certs) == 0 || len(certs[0]) == 0 {
		return errors.New("missing certificate chain")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	userName := certs[0][0].Subject.CommonName

	res, err := auth.authApp.Authorize(userName, method)
	if err != nil {
		return errors.WithStack(err)
	}

	allowed := res
	if !allowed {
		return errors.New("Unauthorized Access")
	}

	return nil
}

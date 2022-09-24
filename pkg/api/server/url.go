package server

import (
	"context"
	"tinyurl/pkg/api/proto"
	"tinyurl/pkg/datamodel"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type tinyUrlServer struct {
	app datamodel.TinyUrlApplication
}

func NewTinyUrlServer(app datamodel.TinyUrlApplication) (proto.TinyUrlServiceServer, error) {
	server := &tinyUrlServer{
		app: app,
	}
	return server, nil
}

func (s *tinyUrlServer) Create(ctx context.Context, r *proto.UrlRequest) (*proto.UrlResponse, error) {
	longUrl := r.GetUrl()
	tinyUrl := s.app.Create(longUrl)
	res := &proto.UrlResponse{
		Url: tinyUrl,
	}
	return res, nil
}

func (s *tinyUrlServer) Fetch(ctx context.Context, r *proto.UrlRequest) (*proto.UrlResponse, error) {
	tinyUrl := r.GetUrl()
	longUrl, err := s.app.Fetch(tinyUrl)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	res := &proto.UrlResponse{
		Url: longUrl,
	}
	return res, nil
}

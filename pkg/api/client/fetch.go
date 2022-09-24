package client

import (
	"context"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"time"
	"tinyurl/pkg/api/proto"

	"google.golang.org/grpc"
)

type fetchCommand struct {
	client   proto.TinyUrlServiceClient
	logEntry *logrus.Entry
}

func NewFetchCommand(client proto.TinyUrlServiceClient, logEntry *logrus.Entry) *fetchCommand {
	return &fetchCommand{
		client:   client,
		logEntry: logEntry,
	}
}

func (c *fetchCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("not enough arguments to run the command")
	}

	if len(args) > 1 {
		return errors.New("too many arguments to run the command")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	tinyUrl := args[0]
	command := proto.UrlRequest{
		Url: tinyUrl,
	}

	res, err := c.client.Fetch(ctx, &command, grpc.WaitForReady(true),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...))
	if err != nil {
		return errors.WithStack(err)
	}

	c.logEntry.Infof("Fetched long url %s for the input tiny url %s", res.GetUrl(), tinyUrl)
	return nil
}

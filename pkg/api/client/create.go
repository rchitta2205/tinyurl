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

type createCommand struct {
	client   proto.TinyUrlServiceClient
	logEntry *logrus.Entry
}

func NewCreateCommand(client proto.TinyUrlServiceClient, logEntry *logrus.Entry) *createCommand {
	return &createCommand{
		client:   client,
		logEntry: logEntry,
	}
}

func (c *createCommand) Execute(args []string) error {
	if len(args) == 0 {
		return errors.New("not enough arguments to run the command")
	}

	if len(args) > 1 {
		return errors.New("too many arguments to run the command")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	longUrl := args[0]
	command := proto.UrlRequest{
		Url: longUrl,
	}

	res, err := c.client.Create(ctx, &command, grpc.WaitForReady(true),
		grpc_retry.WithCodes(grpc_retry.DefaultRetriableCodes...))
	if err != nil {
		return errors.WithStack(err)
	}

	c.logEntry.Infof("Created a tiny url %s for the input long url %s", res.GetUrl(), longUrl)
	return nil
}

package client

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"tinyurl/pkg/api/proto"
	"tinyurl/pkg/config"
)

type Command interface {
	Execute(args []string) error
}

type commandManager struct {
	ctx      context.Context
	cfg      config.Config
	conn     *clientConn
	logEntry *logrus.Entry

	// List of all clients
	tinyUrlClient proto.TinyUrlServiceClient

	// List of all commands
	createCommand *createCommand
	fetchCommand  *fetchCommand
}

type CommandManagerOption func(cm *commandManager)

func NewCommandManager(ctx context.Context, cfg config.Config, logEntry *logrus.Entry, opts ...CommandManagerOption) (*commandManager, error) {
	var err error
	cm := &commandManager{
		ctx:      ctx,
		cfg:      cfg,
		logEntry: logEntry,
	}

	cm.conn = NewClientConn(ctx, cfg)
	err = cm.conn.dial()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, opt := range opts {
		opt(cm)
	}
	return cm, nil
}

// Options pattern to initialize all the required clients
func WithTinyUrlServiceClient() CommandManagerOption {
	return func(cm *commandManager) {
		cm.tinyUrlClient = proto.NewTinyUrlServiceClient(cm.conn.getGrpcConn())
	}
}

func (cm *commandManager) CreateCommand() Command {
	if cm.createCommand == nil {
		cm.createCommand = NewCreateCommand(cm.tinyUrlClient, cm.logEntry)
	}
	return cm.createCommand
}

func (cm *commandManager) FetchCommand() Command {
	if cm.fetchCommand == nil {
		cm.fetchCommand = NewFetchCommand(cm.tinyUrlClient, cm.logEntry)
	}
	return cm.fetchCommand
}

func (cm *commandManager) Release() error {
	return cm.conn.release()
}

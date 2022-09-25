package pkg

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	api "tinyurl/pkg/api/client"
	"tinyurl/pkg/config"
)

type client struct {
	ctx      context.Context
	cfg      config.Config
	args     []string
	cmds     map[string]api.Command
	logEntry *logrus.Entry
	release  func() error
}

func NewClient(ctx context.Context, cfg config.Config, args []string, logEntry *logrus.Entry) *client {
	return &client{
		ctx:      ctx,
		cfg:      cfg,
		args:     args,
		logEntry: logEntry,
	}
}

func (c *client) Register() error {

	var cmdManagerOpts []api.CommandManagerOption

	// Initialize all the clients
	cmdManagerOpts = append(cmdManagerOpts, api.WithTinyUrlServiceClient())

	// Create the command manager
	cmdMgr, err := api.NewCommandManager(c.ctx, c.cfg, c.logEntry, cmdManagerOpts...)
	if err != nil {
		return errors.WithStack(err)
	}
	c.release = cmdMgr.Release

	// Initialize all the used commands in the application
	c.cmds = map[string]api.Command{
		"create": cmdMgr.CreateCommand(),
		"fetch":  cmdMgr.FetchCommand(),
	}

	return nil
}

func (c *client) Run() error {
	if len(c.args) < 1 {
		return errors.New("Not enough arguments to run the command")
	}

	cmd, ok := c.cmds[c.args[0]]
	if ok {
		return cmd.Execute(c.args[1:])
	}

	return errors.New("unknown command: " + c.args[0])
}

func (c *client) Release() error {
	return c.release()
}

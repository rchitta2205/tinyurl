package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"os"
	"tinyurl/pkg"
	"tinyurl/pkg/config"
)

func main() {
	ctx := context.Background()
	cfg := config.NewConfig()
	logEntry := logrus.NewEntry(logrus.New())

	client := pkg.NewClient(ctx, cfg, os.Args[1:], logEntry)

	err := client.Register()
	if err != nil {
		os.Exit(-1)
	}

	err = client.Run()
	if err != nil {
		os.Exit(-1)
	}

	err = client.Release()
	if err != nil {
		os.Exit(-1)
	}
}

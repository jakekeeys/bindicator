package main

import (
	"context"
	"os"

	"github.com/jakekeeys/bindicator/internal/api"
	"github.com/sirupsen/logrus"
)

func main() {
	err := run(context.Background())
	if err != nil {
		logrus.WithError(err).Panic("could not run bindicator")
	}
}

func run(ctx context.Context) error {
	logrus.Info("running bin collection app")

	debug := os.Getenv("DEBUG") == "true"
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug mode enabled")
	}

	return api.NewHTTP(ctx, debug)
}

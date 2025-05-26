/*
Copyright © 2025 Joao Cervino jcervinobarbosa@gmail.com
*/

package state

import (
	"context"
	"os"

	"github.com/sirupsen/logrus"
)

type GlobalState struct {
	Ctx     context.Context
	Logger  *logrus.Logger
	Docker  bool
	OSExit  func(int)
	CfgFile string
}

func NewState(ctx context.Context) *GlobalState {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetLevel(logrus.InfoLevel)

	return &GlobalState{
		Ctx:     ctx,
		Docker:  false,
		Logger:  logger,
		OSExit:  os.Exit,
		CfgFile: "./scannercli.yaml",
	}
}

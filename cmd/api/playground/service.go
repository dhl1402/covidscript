package playground

import (
	"bytes"
	"context"
	"fmt"

	"github.com/dhl1402/covidscript/internal/config"
	"github.com/dhl1402/covidscript/internal/interpreter"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

type Service interface {
	Intepret(context.Context, InterpretRequest) (*InterpretResponse, error)
}

type service struct {
	logger *zap.SugaredLogger
}

func New(logger *zap.SugaredLogger) Service {
	return &service{
		logger: logger,
	}
}

func (s *service) Intepret(ctx context.Context, req InterpretRequest) (*InterpretResponse, error) {
	var errMessage string
	var buf bytes.Buffer
	err := interpreter.Interpret(req.Script, config.Config{Writer: &buf})
	if err != nil {
		errMessage = err.Error()
		sentry.CaptureException(fmt.Errorf("Script: %s\nError: %s", req.Script, errMessage))
	}
	return &InterpretResponse{
		Error:    errMessage,
		Response: []string{buf.String()},
	}, nil
}

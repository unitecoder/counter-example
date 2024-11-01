package main

import (
	"context"
	"fmt"
	"github.com/functionstream/function-stream/clients/gofs"
	"github.com/functionstream/function-stream/common"
	"strconv"
)

func main() {
	err := gofs.NewFSClient().
		Register(gofs.DefaultModule, gofs.Function(counter).AddInitFunc(initCounter)).
		Run()
	if err != nil {
		panic(err)
	}
}

type Config struct {
	ServiceName string `json:"service_name"`
}

type Event struct {
	User string `json:"user"`
}

type Result struct {
	Event
	ServiceName string `json:"serviceName"`
	Count       int    `json:"count"`
}

var config Config

func initCounter(ctx context.Context) error {
	funcCtx := gofs.GetFunctionContext(ctx)
	configMap, err := funcCtx.GetConfig(ctx)
	if err != nil {
		return fmt.Errorf("failed to parse config: %s", err)
	}
	config.ServiceName = configMap["serviceName"]
	return nil
}

func counter(ctx context.Context, e *Event) *Result {
	logger := common.NewDefaultLogger()
	funcCtx := gofs.GetFunctionContext(ctx)
	var count int
	state, err := funcCtx.GetState(ctx, e.User)
	if err != nil {
		err := funcCtx.PutState(ctx, e.User, []byte("0"))
		if err != nil {
			logger.Error(err, "Failed to put state")
			return nil
		}
		logger.Info("State not found, initialized to 0", "user", e.User)
		count = 0
	} else {
		count, err = strconv.Atoi(string(state))
		if err != nil {
			logger.Error(err, "Failed to convert state to integer")
			return nil
		}
	}
	count++
	err = funcCtx.PutState(ctx, e.User, []byte(strconv.Itoa(count)))
	if err != nil {
		logger.Error(err, "Failed to put state")
		return nil
	}
	return &Result{
		Event:       *e,
		ServiceName: config.ServiceName,
		Count:       count,
	}
}

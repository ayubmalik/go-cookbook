package main

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// ***************************************************************************
	// development mode, normall would set this based on some flag or env property
	// ***************************************************************************
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("test message 1")
	logger.Debug("test message 2")
	logger.Warn("test message 3")

	// with an additional log fields
	logger.Error(
		"User failed age validation",
		zapcore.Field{
			Key:     "age",
			Type:    zapcore.Int32Type,
			Integer: 17,
		},
		zapcore.Field{
			Key:    "username",
			Type:   zapcore.StringType,
			String: "some user",
		},
	)

	// same as above but with sugar :)
	sugar := logger.Sugar()
	sugar.Info("test message 1")
	sugar.Infof("test message 1%s", "A")
	sugar.Errorw("User failed age validation",
		"age", 17,
		"username", "some user",
	)

	// ***************************************************************************
	// production mode, normally would set this based on some flag or env property
	// ***************************************************************************
	logger, err = zap.NewProduction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer logger.Sync()
}

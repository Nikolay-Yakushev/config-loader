package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nikolay-Yakushev/config-loader/internal/application"
	logger "github.com/Nikolay-Yakushev/config-loader/pkg/logger"
)
func main(){
	logger, err := logger.New()
	if err != nil{
		log.Fatalf("logger initialization failed: %s", err.Error())
	}
	logger.Sugar().Debug("Logger init compleated")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGKILL, os.Interrupt)
	defer cancel()

	app, err := application.New(context.Background(), logger)
	if err != nil {
		log.Fatalf("App initialization failed: %s", err.Error())
	}
	logger.Sugar().Debugf("App(description=%s) initialization compleated", app.Descr)
	app.Start()
	
	<-ctx.Done()
	stopCtx, stopCancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer stopCancel()
	
	app.Stop(stopCtx)

}
package http

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	ports "github.com/Nikolay-Yakushev/config-loader/internal/ports/driver-ports"
	envcfg "github.com/Nikolay-Yakushev/config-loader/pkg/configs/envcfg"
	filecfg "github.com/Nikolay-Yakushev/config-loader/pkg/configs/filecfg"
)

type Config struct {
	Port int `env:"PORT" envDefault:"3000"`
}

type Adapter struct {
	srv     *http.Server
	log     *zap.Logger
	l       net.Listener
	once    sync.Once
	storage map[string]ports.Config
}

func New(log *zap.Logger) (*Adapter, error) {
	// TODO var storage make(map[string]ports.Config) - does not work 
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("parse server http adapter configuration failed: %w", err)
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		return nil, fmt.Errorf("server start failed: %w", err)
	}
	router := gin.New()

	server := &http.Server{
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// initiate fileconfig
	filecfg, err := filecfg.New()
	if err !=nil {
		return nil, fmt.Errorf("Failed to init storage: %w", err)
	}
	storage := make(map[string]ports.Config)
	storage["file"] = filecfg
	// initiate envconfig
	envcfg, err := envcfg.New()
	if err !=nil {
		return nil, fmt.Errorf("Failed to init storage: %w", err)
	}
	storage["env"] = envcfg
	
	adap := &Adapter{
		srv: server,
		log: log,
		l:   l,
		storage: storage,
	}
	
	initRoutes(adap, router, log)
	return adap, nil
}

func (a *Adapter) Start() error {
	var err error

	go func() {
		err = a.srv.Serve(a.l)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			err = fmt.Errorf("server start failed %w", err)
		}
		err = nil
	}()

	if err != nil {
		return err
	}

	return nil
}

func (a *Adapter) Stop(ctx context.Context) error {
	var err  error
	a.once.Do(func() {
		err = a.srv.Shutdown(ctx)
	})
	a.log.Sugar().Info("Web adapter shut down successefully")
	return err
}

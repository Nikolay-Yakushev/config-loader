package application

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/exp/slices"

	http "github.com/Nikolay-Yakushev/config-loader/internal/adapters/http"
)

type App struct {
	ctx           context.Context
	Log           *zap.Logger
	shutDownFuncs []func(ctx context.Context) error
	Descr         string
}

func New(ctx context.Context, logger *zap.Logger) (*App, error) {

	app := &App{
		ctx:   ctx,
		Log:   logger,
		Descr: "App",
	}

	return app, nil
}

func (a *App) Start() error {
	webAdapter, err := http.New(a.Log)
	if err != nil {
		a.Log.Sugar()
	}
	a.shutDownFuncs = slices.Insert(a.shutDownFuncs, 0, webAdapter.Stop)
	webAdapter.Start()
	return nil
}

func (a *App) Stop(ctx context.Context) error {
	defer a.Log.Sync()

	for _, closeFunc := range a.shutDownFuncs {
		err := closeFunc(a.ctx)
		if err != nil {
			a.Log.Sugar().Error(err)
		}
	}
	a.Log.Sugar().Info("application stopped")
	return nil
}

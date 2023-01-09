package logger

import (
	"fmt"

	"go.uber.org/zap"
)

func New() (logger *zap.Logger, err error) {
	zapCfg := zap.NewProductionConfig()
	logger, err = zapCfg.Build()
	if err != nil {
		err := fmt.Errorf("Cant intiate logger: %w", err)
		return nil, err
	}
	return logger, nil
}

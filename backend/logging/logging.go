package logging

import (
	"go.uber.org/zap"
)

func New(options ...zap.Option) (*zap.Logger, error) {
	return zap.NewProductionConfig().Build(options...)
}

package foundation

import (
	"context"
	"sync"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

const ContextDBKey = "DB"
const ContextLoggerKey = "logger"
const ContextRequestIdKey = "requestId"
const ContextUserClaimKey = "userClaim"
const ContextDomainEventPublisherKey = "domainEventPublisher"

var (
	contextProviderOnce     sync.Once
	contextProviderInstance *contextProvider
)

func ContextProvider() *contextProvider {
	contextProviderOnce.Do(func() {
		contextProviderInstance = &contextProvider{}
	})

	return contextProviderInstance
}

type contextProvider struct {
}

func (contextProvider) GetDB(ctx context.Context) *gorm.DB {
	v := ctx.Value(ContextDBKey)
	if v == nil {
		panic("DB is not exist")
	}
	if db, ok := v.(*gorm.DB); ok {
		return db
	}
	panic("DB is not exist")
}

func (contextProvider) SetDB(ctx context.Context, gormDB *gorm.DB) context.Context {
	return context.WithValue(ctx, ContextDBKey, gormDB)
}

func (contextProvider) SetLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLoggerKey, logger)
}

func (contextProvider) GetLogger(ctx context.Context) *zap.Logger {
	v := ctx.Value(ContextLoggerKey)
	if v == nil {
		panic("Logger is not exist")
	}
	if logger, ok := v.(*zap.Logger); ok {
		return logger
	}
	panic("Logger is not exist")
}

func (contextProvider) SetRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, ContextRequestIdKey, requestId)
}

func (contextProvider) GetRequestId(ctx context.Context) string {
	v := ctx.Value(ContextRequestIdKey)
	if v == nil {
		panic("RequestId is not exist")
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	panic("RequestId is not exist")
}

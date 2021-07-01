package dbrutil

import (
	"context"
	"github.com/gocraft/dbr/v2"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
)

type key int

const (
	dbSessionKey key = 1
	dbTxKey      key = 2
	dbLoggerKey  key = 3
)

func NewContextWithDbSession(ctx context.Context, session *dbr.Session) context.Context {
	return context.WithValue(ctx, dbSessionKey, session)
}

func DbSessionFromContext(ctx context.Context) (*dbr.Session, bool) {
	session, ok := ctx.Value(dbSessionKey).(*dbr.Session)
	return session, ok
}

func NewContextWithDbTx(ctx context.Context, tx *dbr.Tx) context.Context {
	return context.WithValue(ctx, dbTxKey, tx)
}

func DbTxFromContext(ctx context.Context) (*dbr.Tx, bool) {
	tx, ok := ctx.Value(dbTxKey).(*dbr.Tx)
	return tx, ok
}

func NewContextWithLogger(ctx context.Context, logger *LogrusEventReceiver) context.Context {
	return context.WithValue(ctx, dbLoggerKey, logger)
}

func LoggerFromContext(ctx context.Context) *LogrusEventReceiver {
	logger, ok := ctx.Value(dbLoggerKey).(*LogrusEventReceiver)
	if !ok {
		return NewLogrusLogger(ctxlogrus.Extract(ctx))
	}
	return logger
}

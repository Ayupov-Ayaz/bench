package main

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"io"
	"testing"
)

func newZerolog() zerolog.Logger {
	return zerolog.New(io.Discard).With().Timestamp().Logger()
}

func newInfoZerolog() zerolog.Logger {
	return newZerolog().Level(zerolog.InfoLevel)
}

const (
	roomID  = "room_id"
	userID  = "user_id"
	gameID  = "game_id"
	message = "failed to send message"
)

type logKey struct{}

func WithZapLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, logKey{}, logger)
}

func ZapLogger(ctx context.Context) *zap.Logger {
	return ctx.Value(logKey{}).(*zap.Logger)
}

func BenchmarkLoggers_Simple(b *testing.B) {
	err := errors.New("error")

	b.Run("1. Zero.Simple()", func(b *testing.B) {
		logger := newInfoZerolog()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Err(err).Str(roomID, roomID).Str(userID, userID).Str(gameID, gameID).Msg(message)
			}
		})
	})

	b.Run("1. Zap.Simple()", func(b *testing.B) {
		logger := newZapLogger()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Error(message, zap.Error(err),
					zap.String(roomID, roomID), zap.String(userID, userID), zap.String(gameID, gameID))
			}
		})
	})

	b.Run("2. Zero.With()", func(b *testing.B) {
		logger := newInfoZerolog()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				child := logger.With().Str(roomID, roomID).Str(userID, userID).Str(gameID, gameID).Logger()
				child.Error().Err(err).Msg(message)
			}
		})
	})

	b.Run("2. Zap.With()", func(b *testing.B) {
		logger := newZapLogger()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.With(zap.Error(err),
					zap.String(roomID, roomID), zap.String(userID, userID), zap.String(gameID, gameID)).Info(message)
			}
		})
	})

	b.Run("3. Zero.Simple().FromCtx()", func(b *testing.B) {
		log := newInfoZerolog()
		ctx := context.Background()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ctx := log.WithContext(ctx)
				zerolog.Ctx(_ctx).Err(err).Str(roomID, roomID).Str(userID, userID).Str(gameID, gameID).Msg(message)
			}
		})
	})

	b.Run("3. Zap.Simple().FromCtx()", func(b *testing.B) {
		log := newZapLogger()
		ctx := context.Background()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ctx := WithZapLogger(ctx, log)
				ZapLogger(_ctx).Error(message, zap.Error(err),
					zap.String(roomID, roomID), zap.String(userID, userID), zap.String(gameID, gameID))
			}
		})
	})

	b.Run("4. Zero.With().FromCtx()", func(b *testing.B) {
		log := newInfoZerolog()
		ctx := context.Background()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ctx := log.With().Str(roomID, roomID).
					Str(userID, userID).Str(gameID, gameID).Logger().WithContext(ctx)
				zerolog.Ctx(_ctx).Error().Err(err).Msg(message)
			}
		})
	})

	b.Run("4. Zap.With().FromCtx()", func(b *testing.B) {
		log := newZapLogger()
		ctx := context.Background()
		b.ResetTimer()
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_ctx := WithZapLogger(ctx, log.With(zap.String(roomID, roomID),
					zap.String(userID, userID), zap.String(gameID, gameID)))
				ZapLogger(_ctx).Error(message, zap.Error(err))
			}
		})
	})
}

//1._Zero.Simple()-8              		500035126               69.32 ns/op            0 B/op          0 allocs/op
//1._Zap.Simple()-8               		127405725              299.7 ns/op           256 B/op          1 allocs/op
//2._Zero.With()-8                		133400842              241.1 ns/op           512 B/op          1 allocs/op
//2._Zap.With()-8                 		45481430               773.3 ns/op          1538 B/op          6 allocs/op
//3._Zero.Simple().FromCtx()-8          284573256              120.7 ns/op           144 B/op          2 allocs/op
//3._Zap.Simple().FromCtx()-8           132590954              281.0 ns/op           304 B/op          2 allocs/op
//4._Zero.With().FromCtx()-8            100000000              300.8 ns/op           656 B/op          3 allocs/op
//4._Zap.With().FromCtx()-8             52770883               800.7 ns/op          1586 B/op          8 allocs/op

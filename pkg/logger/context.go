/**
 * The Ticket-er Telegram Bot Source Code
 * Copyright (C) 2026 Lev (Leo) Kondukov (aka DiceBarrel, Barrel, Vandi)
 * 
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License.
 * 
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU General Public License for more details.
 * 
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package logger

import (
	"context"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Key string

const (
	ContextLogger Key = "logger-value"
	SessionId     Key = "session-id"
	BotChatId     Key = "bot-chat-id"
	BotMessageId  Key = "bot-message-id"
	BotUpdateId   Key = "bot-update-id"
)

func Context(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ContextLogger, l)
}

func AddId(ctx context.Context, chat int64, msg, update int) (context.Context, string) {
	v := uuid.New()
	return context.WithValue(
		context.WithValue(
			context.WithValue(
				context.WithValue(ctx, BotUpdateId, update), SessionId, v,
			), BotChatId, chat,
		), BotMessageId, msg,
	), v.String()
}

func FromCtx(ctx context.Context) *zap.Logger {
	logger, ok := ctx.Value(ContextLogger).(*zap.Logger)
	if !ok {
		return nil
	}
	return logger
}

func SessionIdFromCtx(ctx context.Context) (uuid.UUID, bool) {
	v, ok := ctx.Value(SessionId).(uuid.UUID)
	return v, ok
}

func BotChatIdFromCtx(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(BotChatId).(int64)
	return v, ok
}

func BotMessageIdFromCtx(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(BotMessageId).(int)
	return v, ok
}

func BotUpdateIdFromCtx(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(BotUpdateId).(int)
	return v, ok
}

func Log(ctx context.Context, lvl zapcore.Level, msg string, fields ...zap.Field) bool {
	logger := FromCtx(ctx)
	if logger == nil {
		return false
	}

	if session, ok := SessionIdFromCtx(ctx); ok {
		fields = append(fields, zap.String(string(SessionId), session.String()))
	}
	if botChatId, ok := BotChatIdFromCtx(ctx); ok {
		fields = append(fields, zap.Int64(string(BotChatId), botChatId))
	}
	if botMessageId, ok := BotMessageIdFromCtx(ctx); ok {
		fields = append(fields, zap.Int(string(BotMessageId), botMessageId))
	}
	if botUpdateId, ok := BotUpdateIdFromCtx(ctx); ok {
		fields = append(fields, zap.Int(string(BotUpdateId), botUpdateId))
	}

	logger.Log(lvl, msg, fields...)
	return true
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DebugLevel, msg, fields...)
}

func Info(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.InfoLevel, msg, fields...)
}

func Warn(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.WarnLevel, msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.ErrorLevel, msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.FatalLevel, msg, fields...)
}

func DPanic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.DPanicLevel, msg, fields...)
}

func Panic(ctx context.Context, msg string, fields ...zap.Field) bool {
	return Log(ctx, zap.PanicLevel, msg, fields...)
}

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
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logPath string) *zap.Logger {
	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		LevelKey:       "level",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		MessageKey:     "message",
		CallerKey:      "caller",
		EncodeDuration: zapcore.StringDurationEncoder,
	})

	jsonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:     "timestamp",
		EncodeTime:  zapcore.ISO8601TimeEncoder,
		LevelKey:    "level",
		EncodeLevel: zapcore.LowercaseLevelEncoder,
		MessageKey:  "message",
	})

	getLogWriter := func(filename string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   filepath.Join(logPath, filename),
			MaxSize:    100, // megabytes
			MaxBackups: 365,
			MaxAge:     365, // days
			Compress:   true,
		})
	}

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zap.DebugLevel),
		zapcore.NewCore(jsonEncoder, getLogWriter("debug/debug.log"), zap.DebugLevel),
		zapcore.NewCore(jsonEncoder, getLogWriter("info/info.log"), zap.InfoLevel),
		zapcore.NewCore(jsonEncoder, getLogWriter("error/error.log"), zap.ErrorLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

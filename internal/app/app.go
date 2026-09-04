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

package app

import (
	"context"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mymmrac/telego"
	"github.com/redis/go-redis/v9"
	"github.com/vandi37/ticket-er/internal/bot"
	"github.com/vandi37/ticket-er/internal/config"
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/pkg/closer"
	"github.com/vandi37/ticket-er/pkg/logger"
	"go.uber.org/zap"
)

func Run(ctx context.Context) {
	logPath := os.Getenv("LOG_PATH")
	l := logger.InitLogger(logPath)
	ctx = context.WithValue(ctx, logger.ContextLogger, l)

	cfg, err := config.LoadConfig()
	if err != nil {
		l.Fatal("error getting configuration", zap.Error(err))
		return
	}
	ctx = context.WithValue(ctx, ctx_vals.ChannelId, cfg.ChannelId)
	split := strings.Split(cfg.Owners, ",")
	owners := make(map[int64]struct{}, len(split))
	for _, v := range split {
		owner, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Fatal("error parsing owners", zap.Error(err))
			return
		}
		owners[owner] = struct{}{}
	}
	ctx = context.WithValue(ctx, ctx_vals.Owners, owners)

	cl := closer.New()

	opt, err := redis.ParseURL(cfg.RedisConnString)
	if err != nil {
		l.Fatal("error parsing redis url", zap.Error(err))
		return
	}

	rdb := redis.NewClient(opt)
	cl.AddNoCtx(rdb.Conn().Close)
	ctx = context.WithValue(ctx, ctx_vals.RedisClient, rdb)

	pool, err := pgxpool.New(ctx, cfg.DbConnString)
	if err != nil {
		l.Fatal("error connecting to postgres", zap.Error(err))
		return
	}
	cl.AddNoError(pool.Close)
	ctx = context.WithValue(ctx, ctx_vals.DbClient, pool)
	if err := pool.Ping(ctx); err != nil {
		l.Fatal("error pining database", zap.Error(err))
	}

	b, err := telego.NewBot(cfg.BotToken)
	if err != nil {
		l.Fatal("error initializing bot", zap.Error(err))
		return
	}
	cl.Add(b.Close)
	me, err := b.GetMe(ctx)
	if err != nil {
		l.Fatal("error getting bot", zap.Error(err))
		return
	}
	logger.Info(ctx, "bot started", zap.String("username", me.Username))
	go bot.Run(ctx, b)

	<-ctx.Done()
	timeout, close := context.WithTimeout(context.Background(), time.Minute)
	defer close()

	err = cl.Close(timeout)
	if err != nil {
		l.Fatal("failed to close", zap.Error(err))
	}
	l.Debug("stopped gracefully")
}

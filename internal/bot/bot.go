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

package bot

import (
	"context"

	th "github.com/mymmrac/telego/telegohandler"

	"github.com/mymmrac/telego"
	"github.com/vandi37/ticket-er/internal/commands"
	"github.com/vandi37/ticket-er/internal/service"
	"github.com/vandi37/ticket-er/internal/usernames"
	"github.com/vandi37/ticket-er/pkg/logger"
	"go.uber.org/zap"
)

func Run(ctx context.Context, b *telego.Bot) {
	updates, err := b.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout: 30,
	})
	if err != nil {
		logger.Error(ctx, "error getting updates", zap.Error(err))
		return
	}
	handler, err := th.NewBotHandler(b, updates)
	if err != nil {
		logger.Error(ctx, "error making handler", zap.Error(err))
		return
	}
	defer func() {
		if err := handler.Stop(); err != nil {
			logger.Error(ctx, "error stopping handler", zap.Error(err))
		}
	}()

	handler.Use(func(c *th.Context, update telego.Update) error {
		cancelCtx, cancel := context.WithCancel(ctx)
		go func() {
			select {
			case <-c.Context().Done():
				cancel()
			case <-cancelCtx.Done():
			}
		}()
		return c.WithContext(cancelCtx).Next(update)
	})
	handler.Use(func(ctx *th.Context, update telego.Update) error {
		if update.Message == nil {
			return nil
		}
		c, _ := logger.AddId(ctx.Context(), update.Message.Chat.ID, update.Message.MessageID, update.UpdateID)
		ctx = ctx.WithContext(c)
		return ctx.Next(update)
	})
	handler.Use(func(ctx *th.Context, update telego.Update) error {
		err := ctx.Next(update)
		if err != nil {
			logger.Error(ctx, "got an error", zap.Error(err))
		}
		return nil
	})
	handler.Use(usernames.UsernamesMiddleware)
	handler.Handle(service.Balance, th.TextMatches(commands.Balance))
	handler.Handle(service.Give, th.TextMatches(commands.Give))
	handler.Handle(service.Take, th.TextMatches(commands.Take))
	handler.Handle(service.Transfer, th.TextMatches(commands.Transfer))
	if err := handler.Start(); err != nil {
		logger.Error(ctx, "error starting handler", zap.Error(err))
	}
}

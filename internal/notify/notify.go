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

package notify

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/internal/models"
	"github.com/vandi37/ticket-er/internal/text"
	"github.com/vandi37/ticket-er/internal/usernames"
	"github.com/vandi37/vanerrors"
)

const ChannelNotFound string = "channel not found"

func Notify(ctx *th.Context, transaction models.Transaction, telegram_sender_id, telegram_receiver_id *int64) error {
	channelId, ok := ctx.Value(ctx_vals.ChannelId).(int64)
	if !ok {
		return vanerrors.Simple(ChannelNotFound)
	}

	var sender *usernames.UserInfo
	var receiver *usernames.UserInfo
	if telegram_sender_id != nil {
		chat, err := ctx.Bot().GetChat(ctx, &telego.GetChatParams{
			ChatID: telego.ChatID{
				ID: *telegram_sender_id,
			},
		})
		if err == nil {
			sender = &usernames.UserInfo{
				ID:       *telegram_sender_id,
				Username: chat.Username,
				Name:     chat.FirstName,
			}
			if len(sender.Name) == 0 {
				sender.Name = chat.Title
			}
		}
	}

	if telegram_receiver_id != nil {
		chat, err := ctx.Bot().GetChat(ctx, &telego.GetChatParams{
			ChatID: telego.ChatID{
				ID: *telegram_receiver_id,
			},
		})
		if err == nil {
			receiver = &usernames.UserInfo{
				ID:       *telegram_receiver_id,
				Username: chat.Username,
				Name:     chat.FirstName,
			}
			if len(receiver.Name) == 0 {
				receiver.Name = chat.Title
			}
		}
	}

	_, err := ctx.Bot().SendMessage(
		ctx,
		tu.Message(telego.ChatID{ID: channelId}, text.TransactionInfo(transaction, sender, receiver)).
			WithParseMode(telego.ModeMarkdownV2).
			WithLinkPreviewOptions(&telego.LinkPreviewOptions{
				IsDisabled: true,
			}),
	)
	return err
}

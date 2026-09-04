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

package reply

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func Reply(ctx *th.Context, message *telego.Message, text string, replyOptions ...ReplyOption) (*telego.Message, error) {
	options := tu.Message(message.Chat.ChatID(), text).
		WithMessageThreadID(message.MessageThreadID).
		WithReplyParameters(&telego.ReplyParameters{
			MessageID:                message.MessageID,
			ChatID:                   message.Chat.ChatID(),
			AllowSendingWithoutReply: true,
		}).WithParseMode(telego.ModeMarkdownV2).
		WithLinkPreviewOptions(&telego.LinkPreviewOptions{
			IsDisabled: true,
		})
	for _, v := range replyOptions {
		v(options)
	}
	return ctx.Bot().SendMessage(ctx, options)
}

type ReplyOption func(message *telego.SendMessageParams)

func WithLinkPreviewOptions(options *telego.LinkPreviewOptions) ReplyOption {
	return func(message *telego.SendMessageParams) {
		message.WithLinkPreviewOptions(options)
	}
}

func WithReplyMarkup(rm telego.ReplyMarkup) ReplyOption {
	return func(message *telego.SendMessageParams) {
		message.WithReplyMarkup(rm)
	}
}

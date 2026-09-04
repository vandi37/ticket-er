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

package service

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/vandi37/ferror"
	"github.com/vandi37/ticket-er/internal/reply"
	"github.com/vandi37/ticket-er/internal/repo"
	"github.com/vandi37/ticket-er/internal/text"
	"github.com/vandi37/ticket-er/internal/usernames"
)

func Balance(ctx *th.Context, update telego.Update) error {
	save := ferror.Save("service.Balance")
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	_, info, ok, err := usernames.GetUserFromMessage(ctx, *update.Message)
	if err != nil {
		return save.New(err)
	}
	if !ok {
		_, err := reply.Reply(ctx, update.Message, text.NotFoundByUsername())
		if err != nil {
			err = save.New(err)
		}
		return err
	}

	_, err = repo.NewUser(ctx, info.ID)
	if err != nil {
		return save.New(err)
	}
	user, err := repo.GetUser(ctx, info.ID, false)
	if err != nil {
		return save.New(err)
	}

	_, err = reply.Reply(ctx, update.Message, text.UserBalance(info, user))
	if err != nil {
		return save.New(err)
	}
	return nil
}

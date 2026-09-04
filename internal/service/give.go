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
	"errors"
	"strconv"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/vandi37/ferror"
	"github.com/vandi37/ticket-er/internal/notify"
	"github.com/vandi37/ticket-er/internal/reply"
	"github.com/vandi37/ticket-er/internal/repo"
	"github.com/vandi37/ticket-er/internal/text"
	"github.com/vandi37/ticket-er/internal/usernames"
	"github.com/vandi37/vanerrors"
)

func Give(ctx *th.Context, update telego.Update) error {
	save := ferror.Save("service.Give")
	if update.Message == nil || update.Message.From == nil {
		return nil
	}
	if ok, err := CheckAdmin(ctx, update.Message); err != nil {
		return save.New(err)
	} else if !ok {
		return nil
	}

	t, info, ok, err := usernames.GetUserFromMessage(ctx, *update.Message)
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
	first, left, _ := strings.Cut(t, "\n")
	_, after, ok := strings.Cut(first, " ")
	if !ok {
		_, err := reply.Reply(ctx, update.Message, text.NotFoundAmount())
		if err != nil {
			err = save.New(err)
		}
		return err
	}
	after = strings.TrimSpace(after)
	amount, err := strconv.ParseInt(strings.ReplaceAll(after, ",", ""), 10, 64)
	if err != nil || amount <= 0 {
		_, err := reply.Reply(ctx, update.Message, text.NotFoundAmount())
		if err != nil {
			err = save.New(err)
		}
		return err
	}

	tx, err := repo.Pay(ctx, 0, info.ID, amount, left, "give", true, false, false)
	if err != nil {
		return save.New(err)
	}
	if err := notify.Notify(ctx, tx, nil, &info.ID); err != nil {
		return save.New(err)
	}

	if _, err := reply.Reply(ctx, update.Message, text.Give(info, amount)); err != nil {
		return save.New(err)
	}
	return nil
}

func Take(ctx *th.Context, update telego.Update) error {
	save := ferror.Save("service.Take")
	if update.Message == nil || update.Message.From == nil {
		return nil
	}
	if ok, err := CheckAdmin(ctx, update.Message); err != nil {
		return save.New(err)
	} else if !ok {
		return nil
	}

	t, info, ok, err := usernames.GetUserFromMessage(ctx, *update.Message)
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
	first, left, _ := strings.Cut(t, "\n")
	_, after, ok := strings.Cut(first, " ")
	if !ok {
		_, err := reply.Reply(ctx, update.Message, text.NotFoundAmount())
		if err != nil {
			err = save.New(err)
		}
		return err
	}
	after = strings.TrimSpace(after)
	amount, err := strconv.ParseInt(strings.ReplaceAll(after, ",", ""), 10, 64)
	if err != nil || amount <= 0 {
		_, err := reply.Reply(ctx, update.Message, text.NotFoundAmount())
		if err != nil {
			err = save.New(err)
		}
		return err
	}

	tx, err := repo.Pay(ctx, 0, info.ID, amount, left, "take", false, true, false)
	if err != nil && errors.Is(err, vanerrors.Simple(repo.INSUFFICIENT_FUNDS)) {
		_, err := reply.Reply(ctx, update.Message, text.InsufficientFunds())
		if err != nil {
			err = save.New(err)
		}
		return err
	}
	if err != nil {
		return save.New(err)
	}
	if err := notify.Notify(ctx, tx, nil, &info.ID); err != nil {
		return save.New(err)
	}

	if _, err := reply.Reply(ctx, update.Message, text.Give(info, amount)); err != nil {
		return save.New(err)
	}
	return nil
}

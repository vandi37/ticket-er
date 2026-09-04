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
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/internal/models"
	"github.com/vandi37/ticket-er/internal/reply"
	"github.com/vandi37/ticket-er/internal/repo"
	"github.com/vandi37/ticket-er/internal/text"
	"github.com/vandi37/vanerrors"
)

const OWNERS_NOT_FOUND = "owners not found"

func CheckAdmin(ctx *th.Context, message *telego.Message) (bool, error) {
	owners, ok := ctx.Value(ctx_vals.Owners).(map[int64]struct{})
	if !ok {
		return false, vanerrors.Simple(OWNERS_NOT_FOUND)
	}
	_, ok = owners[message.From.ID]
	if ok {
		return true, nil
	}
	_, err := repo.NewUser(ctx, message.From.ID)
	if err != nil {
		return false, err
	}
	user, err := repo.GetUser(ctx, message.From.ID, false)
	if err != nil {
		return false, err
	}
	if user.Role == models.Admin {
		return true, nil
	}
	_, err = reply.Reply(ctx, message, text.NotOwner())
	if err != nil {
		return false, err
	}
	return false, nil
}

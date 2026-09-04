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

package usernames

import (
	"context"
	"errors"
	"fmt"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/redis/go-redis/v9"
	"github.com/vandi37/ferror"
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/pkg/extractor"
	"github.com/vandi37/ticket-er/pkg/finder"
	"github.com/vandi37/vanerrors"
)

const (
	REDIS_NOT_FOUND    = "redis not found"
	USERNAME_NOT_FOUND = "username not found"
)

func UsernamesMiddleware(ctx *th.Context, update telego.Update) error {
	save := ferror.Save("usernames.UsernamesMiddleware")
	rdb, ok := ctx.Value(ctx_vals.RedisClient).(*redis.Client)
	if !ok {
		return save.Simple(REDIS_NOT_FOUND)
	}
	users, chats := finder.FindTypes[telego.User, telego.Chat](update)
	for _, v := range users {
		if err := addUsername(ctx, rdb, v.Username, v.ID, v.FirstName); err != nil {
			return save.New(err)
		}
	}
	for _, v := range chats {
		name := v.FirstName
		if len(name) == 0 {
			name = v.Title
		}
		if err := addUsername(ctx, rdb, v.Username, v.ID, name); err != nil {
			return save.New(err)
		}
	}
	return ctx.Next(update)
}

func addUsername(ctx context.Context, rdb *redis.Client, username string, user_id int64, user_name string) error {
	if len(username) == 0 {
		return nil
	}
	id, name, err := getByUsername(ctx, rdb, username)
	if err != nil && err != redis.Nil {
		return err
	}
	if err == nil && id == user_id && name == user_name {
		return nil
	}
	return rdb.Set(ctx, "username:"+username, fmt.Sprintf("%d:%s", user_id, user_name), 0).Err()
}

func getByUsername(ctx context.Context, rdb *redis.Client, username string) (id int64, name string, err error) {
	res, err := rdb.Get(ctx, "username:"+username).Result()
	if err != nil {
		return
	}
	fmt.Sscanf(res, "%d:%s", &id, &name)
	return
}

func GetByUsername(ctx context.Context, username string) (int64, string, error) {
	save := ferror.Save("usernames.GetByUsername")
	rdb, ok := ctx.Value(ctx_vals.RedisClient).(*redis.Client)
	if !ok {
		return 0, "", save.Simple(REDIS_NOT_FOUND)
	}
	id, name, err := getByUsername(ctx, rdb, username)
	if err == redis.Nil {
		return 0, "", save.Simple(USERNAME_NOT_FOUND)
	}
	if err != nil {
		return 0, "", err
	}
	return id, name, nil
}

type UserInfo struct {
	ID       int64
	Username string
	Name     string
}

func GetUserFromMessage(ctx context.Context, message telego.Message) (text string, userInfo UserInfo, ok bool, err error) {
	ok = true
	res := extractor.ExtractFirstMention(message.Text, message.Entities)
	if !res.Found && message.ReplyToMessage != nil && message.ReplyToMessage.From != nil {
		userInfo.ID = message.ReplyToMessage.From.ID
		userInfo.Username = message.ReplyToMessage.From.Username
		userInfo.Name = message.ReplyToMessage.From.FirstName
	} else if !res.Found {
		userInfo.ID = message.From.ID
		userInfo.Username = message.From.Username
		userInfo.Name = message.From.FirstName
	} else if res.User != nil {
		userInfo.ID = res.User.ID
		userInfo.Username = res.User.Username
		userInfo.Name = res.User.FirstName
	} else {
		var id int64
		var name string
		id, name, err = GetByUsername(ctx, res.Username)
		if err != nil && errors.Is(err, vanerrors.Simple(USERNAME_NOT_FOUND)) {
			ok = false
			err = nil
			return
		}
		if err != nil {
			return
		}
		userInfo.ID = id
		userInfo.Username = res.Username
		userInfo.Name = name
	}

	text = res.CleanText
	return
}

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

package repo

import (
	"context"

	"github.com/vandi37/ferror"
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/internal/models"
)

func NewUser(ctx context.Context, telegram_id int64) (bool, error) {
	save := ferror.Save("repo.NewUser")
	db, ok := ctx.Value(ctx_vals.DbClient).(Querier)
	if !ok {
		return false, save.Simple(DB_NOT_FOUND)
	}

	ct, err := db.Exec(ctx, `insert into users (telegram_id) values ($1)
on conflict do nothing`, telegram_id)
	if err != nil {
		return false, save.New(err)
	}
	return ct.RowsAffected() == 1, nil
}

func SetRole(ctx context.Context, telegram_id int64, role models.Role) error {
	save := ferror.Save("repo.SetRole")
	db, ok := ctx.Value(ctx_vals.DbClient).(Querier)
	if !ok {
		return save.Simple(DB_NOT_FOUND)
	}

	ct, err := db.Exec(ctx, `update users set role = $1
where telegram_id=$2 and role <> $1
`, role.GetRole())
	if err != nil {
		return save.New(err)
	}
	if ct.RowsAffected() == 0 {
		return save.Simple(NOTHING_CHANGED)
	}
	return nil
}

func GetUser(ctx context.Context, telegram_id int64, forUpdate bool) (models.User, error) {
	save := ferror.Save("repo.GetUser")
	var user models.User

	db, ok := ctx.Value(ctx_vals.DbClient).(Querier)
	if !ok {
		return user, save.Simple(DB_NOT_FOUND)
	}
	query := `select 
id, 
telegram_id,
balance,
role,
updated_at,
created_at
from users 
where telegram_id = $1
limit 1`
	if forUpdate {
		query += "\nfor update"
	}
	row := db.QueryRow(ctx, query, telegram_id)
	var role string
	if err := row.Scan(&user.Id, &user.TelegramId, &user.Balance, &role, &user.UpdatedAt, &user.CreatedAt); err != nil {
		return user, save.New(err)
	}
	r, ok := models.FromString(role)
	if !ok {
		return user, save.Simple(FAILED_TO_GET_ROLE)
	}
	user.Role = r
	return user, nil
}

func GetFirstPool(ctx context.Context, forUpdate bool) (models.User, error) {
	save := ferror.Save("repo.GetFirstPool")
	var user models.User

	db, ok := ctx.Value(ctx_vals.DbClient).(Querier)
	if !ok {
		return user, save.Simple(DB_NOT_FOUND)
	}
	query := `select 
id, 
balance,
updated_at,
created_at
from users 
where role = $1
limit 1`
	if forUpdate {
		query += "\nfor update"
	}
	row := db.QueryRow(ctx, query, models.Pool.GetRole())
	if err := row.Scan(&user.Id, &user.Balance, &user.UpdatedAt, &user.CreatedAt); err != nil {
		return user, save.New(err)
	}
	user.Role = models.Pool
	return user, nil
}

func GetLeaderboard(ctx context.Context, limit, offset int64) ([]models.User, error) {
	save := ferror.Save("repo.GetLeaderboard")

	db, ok := ctx.Value(ctx_vals.DbClient).(Querier)
	if !ok {
		return nil, save.Simple(DB_NOT_FOUND)
	}
	rows, err := db.Query(ctx, `select
id, 
telegram_id,
balance,
role,
updated_at,
created_at
from users 
order by balance asc limit $1 offset $2
	`, limit, offset)
	if err != nil {
		return nil, save.New(err)
	}
	defer rows.Close()
	users := make([]models.User, 0, limit)
	for rows.Next() {
		var user models.User
		var role string
		if err := rows.Scan(&user.Id, &user.TelegramId, &user.Balance, &role, &user.UpdatedAt, &user.CreatedAt); err != nil {
			return nil, save.New(err)
		}
		r, ok := models.FromString(role)
		if !ok {
			return nil, save.Simple(FAILED_TO_GET_ROLE)
		}
		user.Role = r
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, save.New(err)
	}

	return users, nil
}

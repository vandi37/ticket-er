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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vandi37/ferror"
	"github.com/vandi37/ticket-er/internal/ctx_vals"
	"github.com/vandi37/ticket-er/internal/models"
)

const (
	AMOUNT_NEGATIVE    = "amount negative"
	USER_BLOCKED       = "user blocked"
	INSUFFICIENT_FUNDS = "insufficient funds"
)

func Pay(ctx context.Context, sender_telegram_id, receiver_telegram_id, amount int64, message, _type string, usePoolSender, usePoolReceiver, allowBlocked bool) (models.Transaction, error) {
	save := ferror.Save("repo.Pay")
	var transaction models.Transaction
	if amount < 0 {
		return transaction, save.Simple(AMOUNT_NEGATIVE)
	}

	db, ok := ctx.Value(ctx_vals.DbClient).(*pgxpool.Pool)
	if !ok {
		return transaction, save.Simple(DB_NOT_FOUND)
	}
	tx, err := db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return transaction, save.New(err)
	}
	defer tx.Rollback(ctx)
	ctx = context.WithValue(ctx, ctx_vals.DbClient, tx)

	var sender models.User
	if usePoolSender {
		sender, err = GetFirstPool(ctx, true)
	} else {
		sender, err = GetUser(ctx, sender_telegram_id, true)
	}
	if err != nil {
		return transaction, save.New(err)
	}

	var receiver models.User
	if usePoolReceiver {
		receiver, err = GetFirstPool(ctx, false)
	} else {
		receiver, err = GetUser(ctx, receiver_telegram_id, false)
	}
	if err != nil {
		return transaction, save.New(err)
	}

	if sender.Role == models.Blocked && !allowBlocked {
		return transaction, save.Simple(USER_BLOCKED)
	}

	if sender.Balance < amount && sender.Role != models.Pool {
		return transaction, save.Simple(INSUFFICIENT_FUNDS)
	}

	_, err = tx.Exec(ctx, `update users set balance = balance - $1 where id = $2`, amount, sender.Id)
	if err != nil {
		return transaction, save.New(err)
	}
	_, err = tx.Exec(ctx, `update users set balance = balance + $1 where id = $2`, amount, receiver.Id)
	if err != nil {
		return transaction, save.New(err)
	}
	if err := tx.QueryRow(ctx, `insert into transactions
(sender_id, receiver_id, amount, message, type)
values ($1, $2, $3, $4, $5)
returning id, sender_id, receiver_id, amount, message, type, updated_at, created_at`,
		sender.Id, receiver.Id, amount, message, _type).
		Scan(
			&transaction.Id,
			&transaction.SenderId,
			&transaction.ReceiverId,
			&transaction.Amount,
			&transaction.Message,
			&transaction.Type,
			&transaction.UpdatedAt,
			&transaction.CreatedAt,
		); err != nil {
		return transaction, save.New(err)
	}

	if err := tx.Commit(ctx); err != nil {
		return transaction, save.New(err)
	}
	return transaction, nil
}

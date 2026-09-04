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

package text

import (
	"fmt"
	"strings"

	"github.com/vandi37/ticket-er/internal/models"
	"github.com/vandi37/ticket-er/internal/usernames"
	"github.com/vandi37/ticket-er/pkg/rutime"
	"github.com/vandi37/ticket-er/pkg/tickets"
)

func NotFoundByUsername() string {
	return `❌ Бот не нашёл упомянутого пользователя`
}
func NotFoundAmount() string {
	return "❌ Нужно указать сумму"
}

func NotOwner() string {
	return "❌ Эта команда доступна только админам"
}
func InsufficientFunds() string {
	return "❌ Недостаточно билетиков"
}

func MakeMention(info usernames.UserInfo) string {
	var link string
	if len(info.Username) > 0 {
		link = "t.me/" + info.Username
	} else {
		link = fmt.Sprintf("tg://user?id=%d", info.ID)
	}
	return fmt.Sprintf("[%s](%s)", Escape(info.Name), link)
}

func UserBalance(info usernames.UserInfo, user models.User) string {
	return fmt.Sprintf("👤 %s %s\n\n🎫 %s\n\n📅 Создан %s",
		MakeMention(info),
		RuRole(user.Role),
		tickets.ToRu(user.Balance),
		rutime.Time(user.CreatedAt).Before(),
	)
}

func RuRole(role models.Role) string {
	switch role {
	case models.Admin:
		return "(админ)"
	case models.Blocked:
		return "(заблокированный)"
	case models.Pool:
		return "(системный)"
	default:
		return ""
	}
}

func TransactionInfo(tx models.Transaction, senderInfo *usernames.UserInfo, receiverInfo *usernames.UserInfo) string {
	var sender string
	if senderInfo != nil {
		sender = MakeMention(*senderInfo)
	} else {
		sender = fmt.Sprintf("`%d`", tx.SenderId)
	}
	var receiver string
	if receiverInfo != nil {
		receiver = MakeMention(*receiverInfo)
	} else {
		receiver = fmt.Sprintf("`%d`", tx.ReceiverId)
	}
	return fmt.Sprintf(
		"💸 Транзакция №`%d`\n\nТип:%s\nОтправитель: %s\nПолучатель: %s\nСумма: %s\nСообщение: %s\nСоздана: %s",
		tx.Id,
		tx.Type,
		sender,
		receiver,
		tickets.ToRu(tx.Amount),
		Escape(tx.Message),
		rutime.Time(tx.CreatedAt).String(),
	)
}
func Give(info usernames.UserInfo, amount int64) string {
	return fmt.Sprintf("✅ Выдано %s %s", tickets.ToRu(amount), MakeMention(info))
}

func Take(info usernames.UserInfo, amount int64) string {
	return fmt.Sprintf("✅ Забрано %s у %s", tickets.ToRu(amount), MakeMention(info))
}
func Transfer(info usernames.UserInfo, amount int64) string {
	return fmt.Sprintf("✅ Передано %s %s", tickets.ToRu(amount), MakeMention(info))
}

func Escape(text string) string {
	return strings.NewReplacer(
		"_", "\\_", "*", "\\*", "[", "\\[", "]", "\\]", "(",
		"\\(", ")", "\\)", "~", "\\~", "`", "\\`", ">", "\\>",
		"#", "\\#", "+", "\\+", "-", "\\-", "=", "\\=", "|",
		"\\|", "{", "\\{", "}", "\\}", ".", "\\.", "!", "\\!",
	).Replace(text)
}

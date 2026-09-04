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

package models

import "time"

type Role struct {
	role string
}

var (
	Pool     = Role{role: "pool"}
	Blocked  = Role{role: "blocked"}
	UserRole = Role{role: "user"}
	Admin    = Role{role: "admin"}
)

func (r Role) GetRole() string { return r.role }
func FromString(s string) (r Role, ok bool) {
	ok = true
	switch s {
	case Pool.role:
		r = Pool
	case Blocked.role:
		r = Blocked
	case UserRole.role:
		r = UserRole
	case Admin.role:
		r = Admin
	default:
		ok = false
	}
	return
}

type User struct {
	Id         int64
	TelegramId int64
	Balance    int64
	Role       Role
	UpdatedAt  time.Time
	CreatedAt  time.Time
}

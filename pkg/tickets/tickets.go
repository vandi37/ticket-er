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

package tickets

import (
	"fmt"
	"strconv"
)

func ToRu(amount int64) string {
	shortenedAmount := FormatWithCommas(amount)
	if (amount/10)%10 == 1 {
		return fmt.Sprintf("`%s` билетиков", shortenedAmount)
	} else if amount%10 == 1 {
		return fmt.Sprintf("`%s` билетик", shortenedAmount)
	} else if amount%10 == 2 || amount%10 == 3 || amount%10 == 4 {
		return fmt.Sprintf("`%s` билетика", shortenedAmount)
	} else {
		return fmt.Sprintf("`%s` билетиков", shortenedAmount)
	}

}

func FormatWithCommas(n int64) string {
	s := strconv.FormatInt(n, 10)
	if n < 1000 && n > -1000 {
		return s
	}
	abs := s
	if n < 0 {
		abs = s[1:]
	}
	res := ""
	for i, l := 0, len(abs); i < l; i++ {
		if i > 0 && (l-i)%3 == 0 {
			res += ","
		}
		res += string(abs[i])
	}
	if n < 0 {
		return "-" + res
	}
	return res
}

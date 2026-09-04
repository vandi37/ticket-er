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

package rutime

import "time"

type index uint8

const (
	empty index = iota
	one
	twoToFour
	other
)

func calc(need int) index {
	first := digitInt(need, 1)
	second := digitInt(need, 2)
	switch {
	case need == 0:
		return empty
	case second == 1:
		return other
	case first == 1:
		return one
	case first < 5 && first > 1:
		return twoToFour
	default:
		return other
	}
}

//func second(second int) string {
//	index := calc(second)
//	switch index {
//	case empty:
//		return ""
//	case one:
//		return "секунда"
//	case twoToFour:
//		return "секунды"
//	default:
//		return "секунд"
//	}
//}

func minute(minute int) string {
	index := calc(minute)
	switch index {
	case empty:
		return ""
	case one:
		return "минута"
	case twoToFour:
		return "минуты"
	default:
		return "минут"
	}
}

func hour(hour int) string {
	index := calc(hour)
	switch index {
	case empty:
		return ""
	case one:
		return "час"
	case twoToFour:
		return "часа"
	default:
		return "часов"
	}
}

func day(day int) string {
	index := calc(day)
	switch index {
	case empty:
		return ""
	case one:
		return "день"
	case twoToFour:
		return "дня"
	default:
		return "дней"
	}
}

func month(month int) string {
	index := calc(month)
	switch index {
	case empty:
		return ""
	case one:
		return "месяц"
	case twoToFour:
		return "месяца"
	default:
		return "месяцев"
	}
}

func year(year int) string {
	index := calc(year)
	switch index {
	case empty:
		return ""
	case one:
		return "год"
	case twoToFour:
		return "года"
	default:
		return "лет"
	}
}

func (m Month) String() string {
	switch time.Month(m) {
	case time.January:
		return "января"
	case time.February:
		return "февраля"
	case time.March:
		return "марта"
	case time.April:
		return "апреля"
	case time.May:
		return "мая"
	case time.June:
		return "июня"
	case time.July:
		return "июля"
	case time.August:
		return "августа"
	case time.September:
		return "сентября"
	case time.October:
		return "октября"
	case time.November:
		return "ноября"
	case time.December:
		return "декабря"
	default:
		return time.Month(m).String()
	}
}

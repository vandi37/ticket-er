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

import (
	"fmt"
	"strings"
	"time"
)

type Month time.Month

type Time time.Time

func (t Time) String() string {
	event := time.Time(t)
	year := fmt.Sprintf("в %d году", event.Year())
	month := Month(event.Month()).String()
	monthDay := event.Day()
	hour := event.Hour()
	minute := event.Minute()
	return fmt.Sprintf("%s, %d %s в %d:%d по %s", year, monthDay, month, hour, minute, event.Location().String())
}

func (t Time) compare() (string, int) {
	event := time.Time(t)
	now := time.Now()
	returns := -1
	if now.Before(event) {
		returns = 1
		event, now = now, event
	}

	years := now.Year() - event.Year()
	months := int(now.Month()) - int(event.Month())
	days := now.Day() - event.Day()
	hours := now.Hour() - event.Hour()
	minutes := now.Minute() - event.Minute()

	if minutes < 0 {
		hours--
		minutes += 60
	}
	if hours < 0 {
		days--
		hours += 24
	}
	if days < 0 {
		months--
		lastDayOfMonth := time.Date(event.Year(), event.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		days += lastDayOfMonth
	}
	if months < 0 {
		years--
		months += 12
	}

	var parts []string

	addPart := func(value int, fn func(int) string) {
		if value > 0 {
			parts = append(parts, fmt.Sprintf("%d %s", value, fn(value)))
		}
	}

	addPart(years, year)
	addPart(months, month)
	addPart(days, day)
	addPart(hours, hour)
	addPart(minutes, minute)

	if len(parts) == 0 {
		return "", 0
	}

	return strings.Join(parts, " "), returns
}

func (t Time) Before() string {
	str, cmp := t.compare()
	if cmp == 0 {
		return "только что"
	} else if cmp > 0 {
		return fmt.Sprintf("%s назад", str)
	} else {
		return fmt.Sprintf("через %s", str)
	}
}

func (t Time) After() string {
	str, cmp := t.compare()
	if cmp == 0 {
		return "немного"
	} else {
		return str
	}
}

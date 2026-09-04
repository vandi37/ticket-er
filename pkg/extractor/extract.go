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

package extractor

import (
	"unicode/utf16"

	"github.com/mymmrac/telego"
)

type ExtractionResult struct {
	Username  string
	User      *telego.User
	CleanText string
	Found     bool
}

func ExtractFirstMention(text string, entities []telego.MessageEntity) ExtractionResult {
	utf16Text := utf16.Encode([]rune(text))
	for _, entity := range entities {
		if entity.Type == "mention" || entity.Type == "text_mention" {
			start := entity.Offset
			end := entity.Offset + entity.Length
			if start < 0 || end > len(utf16Text) {
				continue
			}
			mentionRaw := string(utf16.Decode(utf16Text[start:end]))
			cleanText := string(utf16.Decode(utf16Text[:start])) + string(utf16.Decode(utf16Text[end:]))
			result := ExtractionResult{
				CleanText: cleanText,
				Found:     true,
			}
			switch entity.Type {
			case "mention":
				if len(mentionRaw) > 0 && mentionRaw[0] == '@' {
					result.Username = mentionRaw[1:]
				} else {
					result.Username = mentionRaw
				}
			case "text_mention":
				result.User = entity.User
			}

			return result
		}
	}

	return ExtractionResult{CleanText: text, Found: false}
}

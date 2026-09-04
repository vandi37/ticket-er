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

package config

import "github.com/goloop/env"

type Config struct {
	BotToken        string `env:"BOT_TOKEN"`
	Owners          string `env:"OWNERS"`
	DbConnString    string `env:"DB_CONN_STRING"`
	RedisConnString string `env:"REDIS_CONN_STRING"`
	ChannelId       int64  `env:"CHANNEL_ID"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	if err := env.Unmarshal("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

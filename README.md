![Golang](https://img.shields.io/badge/-00ADD8?style=for-the-badge&logo=go&logoColor=white&logoSize=auto)
![Telego](https://img.shields.io/badge/Telego-24A1DE?style=for-the-badge&logo=telegram&logoColor=white&logoSize=auto)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-0064a5?style=for-the-badge&logo=postgresql&logoColor=white&logoSize=auto)
![Goose](https://img.shields.io/badge/Goose-000000?style=for-the-badge&logo=go&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Docker Compose](https://img.shields.io/badge/Docker_Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-d82c20?style=for-the-badge&logo=redis&logoColor=white&logoSize=auto)

# Ticket-er bot

Ticketer is a telegram bot that manages a closed economy with a currency called "tickets". 
Currently the interface is only in russian language (because of who this bot was made for)

## Commands

- /start: start the bot
- /help: get bot commands
- кошелёк (user mention): get user balance
- начислить (amount) (user mention): give user tickets (only for admins)
- забрать (amount) (user mention): take tickets from user (only for admins)
- передать (amount) (user mention): transfer tickets to another user

## Deployment

### Download the project
```shell
git clone https://github.com/vandi37/ticket-er.git
```

### Go to project directory
```shell
cd ticket-er
```

### Make .env
```shell
cp example.env .env
```
(edit the .env file to add your configuration)
- BOT_TOKEN=the token of your bot, that you can get in [bot father](https://t.me/BotFather) 
- OWNERS=the owners' telegram ids (you can get them in various telegram bots) format: 1,2,3,4 if only one: 1
- CHANNEL_ID=the channel where transactions are published to

### Run in docker
```shell
docker compose up --build
```
or
```shell
docker-compose up --build
```

## Tecnology stack
- Fully written in [Golang](https://go.dev/)
- Telegram bot api package: [mymmrac/telego](https://pkg.go.dev/github.com/mymmrac/telego)
- Database [PostgreSQL](https://www.postgresql.org/)
- Database package [jackc/pgx/v5](https://pkg.go.dev/github.com/jackc/pgx/v5)
- Username chashing in [Redis](https://redis.io/)
- Redis package [redis/go-redis/v9](https://pkg.go.dev/mod/github.com/redis/go-redis/v9]
- Migrations [goose](https://github.com/pressly/goose)
- You can find all go packages used in [go.mod](go.mod)

## Other info

- [GPL-3.0 LICENSE](LICENSE)
- [Contribute to the project](CONTRIBUTING.md)
- [Future plans](TODO.md)

# Days Without Shit Bot

A Telegram bot that tracks days since the last reset for each chat/channel and sends daily sticker updates at 12:00.

## Features

- Track days since last reset for each chat separately
- Daily sticker notifications at 12:00
- Different stickers for 0, 1, 2, and 3 days
- Simple commands: `/start`, `/reset`, `/days`

## Setup

### 1. Create a Telegram Bot

1. Open Telegram and search for [@BotFather](https://t.me/botfather)
2. Send `/newbot` command
3. Follow the instructions to create your bot
4. Save the bot token (it looks like `123456789:ABCdefGHIjklmNOPQrstUVwxyz`)

### 2. Get Sticker IDs

You need 4 different stickers for different day counts. Here's how to get their IDs:

1. Create your sticker pack or choose existing stickers
2. Send each sticker to [@@RawDataBot](https://t.me/RawDataBot)
3. The bot will reply with sticker information including the File ID
4. Copy the File IDs and replace them in `main.go`:

```go
var stickers = map[int]string{
    0: "YOUR_STICKER_ID_FOR_0_DAYS",
    1: "YOUR_STICKER_ID_FOR_1_DAY",
    2: "YOUR_STICKER_ID_FOR_2_DAYS",
    3: "YOUR_STICKER_ID_FOR_3_PLUS_DAYS",
}
```

### 3. Build and Run

1. Make sure you have Go installed (version 1.21 or later)
2. Clone this repository
3. Build the bot:
   ```bash
   # For your current platform
   make

   # For Linux AMD64
   make build-amd64
   ```
4. Run the bot:
   ```bash
   export TELEGRAM_BOT_TOKEN=YOUR_BOT_TOKEN
   ./days-without-shit  # or ./days-without-shit-amd64 for Linux AMD64 build
   ```

## Usage

1. Add the bot to your group chat or start a private chat
2. Send `/shit_start` to get started
3. Use `/shit_reset` to reset the counter
4. Use `/shit_days` to check how many days have passed
5. The bot will automatically send stickers at 12:00 each day based on the day count:
   - 0 days: First sticker
   - 1 day: Second sticker
   - 2 days: Third sticker
   - 3+ days: Fourth sticker

## Running as a Service

To keep the bot running continuously, you can set it up as a systemd service:

1. Create a service file:
   ```bash
   sudo nano /etc/systemd/system/days-without-shit.service
   ```

2. Add the following content:
   ```ini
   [Unit]
   Description=Days Without Shit Bot
   After=network.target

   [Service]
   ExecStart=/path/to/days-without-shit
   Environment=TELEGRAM_BOT_TOKEN=your_bot_token
   Restart=always
   User=your_username

   [Install]
   WantedBy=multi-user.target
   ```

3. Enable and start the service:
   ```bash
   sudo systemctl enable days-without-shit
   sudo systemctl start days-without-shit
   ```

4. Check status:
   ```bash
   sudo systemctl status days-without-shit
   ```

## State Storage

The bot stores its state in `state.json` in the same directory as the executable. Make sure the directory is writable by the bot process.

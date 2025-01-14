package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ChatState struct {
	LastReset    time.Time `json:"last_reset"`
	LastNotified time.Time `json:"last_notified"`
}

type State struct {
	Chats map[int64]ChatState `json:"chats"`
}

// Sticker IDs for different day counts
var stickers = map[int]string{
	0: "CAACAgIAAxkBAAEw0yBnhhAcJMHwjwOnlBgJevia-vEVeQACFhgAAq_XuEpsPPAf-bYb9TYE", // replace with actual sticker IDs
	1: "CAACAgIAAxkBAAEw0yJnhhA5ePC5ey_Ngzt51qxm0aP7eAACDB8AAgLIsEp4xyhn1z4VZDYE",
	2: "CAACAgIAAxkBAAEw0yRnhhBJT-lX2PUr1ls00RDfkHG_HgACUx8AAhpRsUoXXxbAzeij7jYE",
	3: "CAACAgIAAxkBAAEw0yZnhhBWUoacTPnbSiFGiv8M7LB3NwACMBYAAugUuEoMYV7rmEKTrTYE",
}

func loadState() State {
	data, err := os.ReadFile("state.json")
	if err != nil {
		return State{
			Chats: make(map[int64]ChatState),
		}
	}

	var state State
	if err := json.Unmarshal(data, &state); err != nil {
		log.Printf("Error unmarshaling state: %v", err)
		return State{
			Chats: make(map[int64]ChatState),
		}
	}

	if state.Chats == nil {
		state.Chats = make(map[int64]ChatState)
	}

	return state
}

func saveState(state State) error {
	data, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("error marshaling state: %w", err)
	}

	return os.WriteFile("state.json", data, 0644)
}

func sendDailyStickers(bot *tgbotapi.BotAPI, state *State) {
	now := time.Now()
	if now.Hour() != 12 || now.Minute() != 0 {
		return
	}

	for chatID, chatState := range state.Chats {
		// Check if we already notified today
		if chatState.LastNotified.Year() == now.Year() &&
			chatState.LastNotified.Month() == now.Month() &&
			chatState.LastNotified.Day() == now.Day() {
			continue
		}

		days := int(now.Sub(chatState.LastReset).Hours() / 24)

		if days > 3 {
			continue
		}

		if stickerID, ok := stickers[days]; ok {
			msg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
			if _, err := bot.Send(msg); err != nil {
				log.Printf("Error sending sticker to chat %d: %v", chatID, err)
				continue
			}

			// Update last notified time
			chatState.LastNotified = now
			state.Chats[chatID] = chatState
			if err := saveState(*state); err != nil {
				log.Printf("Error saving state after notification: %v", err)
			}
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	state := loadState()

	// Create a ticker for checking time every minute
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Start a goroutine for sending daily stickers
	go func() {
		for range ticker.C {
			sendDailyStickers(bot, &state)
		}
	}()

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		msg := tgbotapi.NewMessage(chatID, "")

		switch update.Message.Command() {
		case "shit_start":
			msg.Text = "Welcome! Use /shit_reset to reset the counter or /shit_days to see how many days have passed since the last reset."
		case "shit_reset":
			state.Chats[chatID] = ChatState{
				LastReset:    time.Now(),
				LastNotified: time.Time{}, // Reset last notification time
			}
			if err := saveState(state); err != nil {
				log.Printf("Error saving state: %v", err)
				msg.Text = "Error saving state"
			} else {
				msg.Text = "Counter has been reset!"
			}
		case "shit_days":
			chatState, exists := state.Chats[chatID]
			if !exists {
				msg.Text = "Counter hasn't been started yet. Use /shit_reset to start counting!"
			} else {
				days := int(time.Since(chatState.LastReset).Hours() / 24)

				if days > 3 {
					msg.Text = fmt.Sprintf("Дней без срача: %d", days)
					if _, err := bot.Send(msg); err != nil {
						log.Printf("Error sending message: %v", err)
					}
					continue
				}

				// Send sticker first
				if stickerID, ok := stickers[days]; ok {
					stickerMsg := tgbotapi.NewSticker(chatID, tgbotapi.FileID(stickerID))
					if _, err := bot.Send(stickerMsg); err != nil {
						log.Printf("Error sending sticker: %v", err)
					}
				}
			}
		default:
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}
}

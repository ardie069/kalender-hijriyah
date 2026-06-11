package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Notifier struct {
	bot *tgbotapi.BotAPI
}

func New(token string) (*Notifier, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Notifier{
		bot: bot,
	}, nil
}

func (n *Notifier) Send(
	ctx context.Context,
	chatID int64,
	message string,
) error {
	msg := tgbotapi.NewMessage(chatID, message)

	_, err := n.bot.Send(msg)

	return err
}

func (n *Notifier) StartListener(ctx context.Context, handler func(chatID int64, text string)) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := n.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message updates
				continue
			}

			// Call the handler
			handler(update.Message.Chat.ID, update.Message.Text)
		}
	}
}

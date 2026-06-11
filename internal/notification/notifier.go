package notification

import "context"

type Notifier interface {
	Send(
		ctx context.Context,
		chatID int64,
		message string,
	) error
}

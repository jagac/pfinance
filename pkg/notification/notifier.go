package notification

import "context"

type Notifier interface {
	Send(ctx context.Context, to, subject, body string) error
}

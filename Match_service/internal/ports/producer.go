package ports

import "context"

type Producer interface {
	Publish(ctx context.Context, topic string, event interface{}) error
}

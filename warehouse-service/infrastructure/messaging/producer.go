package messaging

import "context"

type Producer interface {
	Publish(ctx context.Context, key, value []byte) error
	Close() error
}

type Message struct {
	Key       []byte
	Value     []byte
	Topic     string
	Partition int
	Offset    int64
}

type Consumer interface {
	ReadMessage(ctx context.Context) (Message, error)
	Commit(ctx context.Context, msg Message) error
	Close() error
}

package domain

import "context"

type Event struct {
	Destination string
	Data        any
}

type EventSender interface {
	Broadcast(ctx context.Context, event Event) error
}

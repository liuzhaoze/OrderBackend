package domain

import "context"

type Event struct {
	Destination string
	Data        any
}

type EventSender interface {
	Direct(ctx context.Context, event Event) error
}

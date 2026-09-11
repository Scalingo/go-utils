package channelsandfunctions

import (
	"context"
	"time"
)

type Event struct{}

type Service interface {
	Events() <-chan *Event
	Subscribe(callback func(ctx context.Context, timer *time.Timer, id string) error) (chan<- Event, func() error)
}
